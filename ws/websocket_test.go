package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/apus-run/van/ws/mocks"
)

func TestConnection_Bind_Success(t *testing.T) {
	upgrader := websocket.Upgrader{}

	tests := []struct {
		name         string
		inputMessage []byte
		expectedData interface{}
	}{
		{
			name:         "Bind to string",
			inputMessage: []byte("Hello, WebSocket"),
			expectedData: "Hello, WebSocket",
		},
		{
			name:         "Bind to JSON struct",
			inputMessage: []byte(`{"key":"value"}`),
			expectedData: map[string]interface{}{"key": "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				conn, err := upgrader.Upgrade(w, r, nil)
				assert.NoError(t, err)
				defer conn.Close()

				wsConn := &Connection{Conn: conn}

				var data interface{}
				switch tt.expectedData.(type) {
				case string:
					data = new(string)
				default:
					data = &map[string]interface{}{}
				}

				err = wsConn.Bind(data)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedData, dereference(data))
			}))
			defer server.Close()

			url := "ws" + server.URL[len("http"):] + "/ws"
			dialer := websocket.DefaultDialer
			conn, resp, err := dialer.Dial(url, nil)
			require.NoError(t, err)

			defer conn.Close()
			defer resp.Body.Close()

			err = conn.WriteMessage(websocket.TextMessage, tt.inputMessage)
			require.NoError(t, err)
		})

		// waiting for previous connection to close and test for new testcase.
		time.Sleep(500 * time.Millisecond)
	}
}

func TestNewWSUpgrader_WithOptions(t *testing.T) {
	errorHandler := func(_ http.ResponseWriter, _ *http.Request, _ int, _ error) {}

	checkOrigin := func(_ *http.Request) bool {
		return true
	}

	options := []Options{
		WithReadBufferSize(1024),
		WithWriteBufferSize(1024),
		WithHandshakeTimeout(500 * time.Millisecond),
		WithSubprotocols("protocol1", "protocol2"),
		WithCompression(),
		WithError(errorHandler),
		WithCheckOrigin(checkOrigin),
	}

	upgrader := NewWSUpgrader(options...)
	actualUpgrader := upgrader.Upgrader.(*websocket.Upgrader)

	assert.Equal(t, 1024, actualUpgrader.ReadBufferSize)
	assert.Equal(t, 1024, actualUpgrader.WriteBufferSize)
	assert.Equal(t, 500*time.Millisecond, actualUpgrader.HandshakeTimeout)
	assert.Equal(t, []string{"protocol1", "protocol2"}, actualUpgrader.Subprotocols)
	assert.True(t, actualUpgrader.EnableCompression)
	assert.NotNil(t, actualUpgrader.Error)
	assert.NotNil(t, actualUpgrader.CheckOrigin)
}

func Test_Upgrade(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpgrader := mocks.NewMockUpgrader(ctrl)

	expectedConn := &websocket.Conn{}
	mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedConn, nil)

	wsUpgrader := WSUpgrader{Upgrader: mockUpgrader}

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, "/", http.NoBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	conn, err := wsUpgrader.Upgrade(w, req, nil)
	require.NoError(t, err)

	assert.Equal(t, expectedConn, conn)
}

func Test_UnimplementedMethods(t *testing.T) {
	conn := &Connection{}

	assert.Empty(t, conn.Param("test"))
	assert.Empty(t, conn.PathParam("test"))
	assert.Empty(t, conn.HostName())
	assert.NotNil(t, conn.Context())
	assert.Nil(t, conn.Params("test"))
}

func dereference(v interface{}) interface{} {
	switch v := v.(type) {
	case *string:
		return *v
	case *map[string]interface{}:
		return *v
	default:
		return v
	}
}

func Test_NewWServer(t *testing.T) {
	r := gin.Default()
	wsManager := New()

	// WebSocket 路由
	r.GET("/ws", func(c *gin.Context) {
		conn, err := wsManager.WebSocketUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			t.Logf("Error during connection upgrade: %v", err)
			c.String(http.StatusInternalServerError, "Could not upgrade connection")
			return
		}

		connID := "some_unique_id"
		// Add the connection to the hub
		connection := &Connection{Conn: conn}
		wsManager.AddWebsocketConnection(connID, connection)

		defer wsManager.CloseConnection(connID)

		for {
			var msg string
			err := connection.Bind(&msg)
			if err != nil {
				t.Logf("Error reading message: %v", err)
				break
			}
			t.Logf("Received message: %s", msg)

			// 发送响应
			message, err := serializeMessage(msg)
			if err != nil {
				t.Logf("Error serializing response: %v", err)
				continue
			}
			err = connection.WriteMessage(TextMessage, message)
			if err != nil {
				t.Logf("Error writing response: %v", err)
				break
			}
		}
	})

	t.Log("WebSocket server started at :8080")
	if err := r.Run(":8080"); err != nil {
		t.Fatal("ListenAndServe:", err)
	}
}

var ErrMarshalingResponse = errors.New("error marshaling response")

func serializeMessage(response any) ([]byte, error) {
	var (
		message []byte
		err     error
	)

	switch v := response.(type) {
	case string:
		message = []byte(v)
	case []byte:
		message = v
	default:
		message, err = json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrMarshalingResponse, err)
		}
	}

	return message, nil
}

func Test_WsCline(t *testing.T) {
	// 连接到 WebSocket 服务器
	url := "ws://localhost:8080/ws"
	dialer := websocket.DefaultDialer
	conn, resp, err := dialer.Dial(url, nil)
	if err != nil {
		t.Fatal("Dial error:", err)
	}
	defer conn.Close()
	defer resp.Body.Close()

	// 发送消息
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, Server!"))
	if err != nil {
		t.Log("Write error:", err)
		return
	}

	// 接收消息
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Log("Read error:", err)
			break
		}
		t.Logf("Received message: %s", msg)

		// 等待一段时间再发送下一条消息
		time.Sleep(1 * time.Second)
	}
}
