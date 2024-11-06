package ws

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WSKey defines the key type for WSConnectionKey.
type WSKey string

// WSConnectionKey is a key constant that stores the connection id in the request context.
const WSConnectionKey WSKey = "ws-connection-key"

// Connection is a wrapper for gorilla websocket connection.
type Connection struct {
	*websocket.Conn
}

type WebSocketFunc func(*Connection) error

// ErrorConnection is the connection error that occurs when webscoket connection cannot be established.
var ErrorConnection = errors.New("couldn't establish connection to web socket")

// The message types are defined in RFC 6455, section 11.8.
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1
)

type WSUpgrader struct {
	Upgrader Upgrader
}

// NewWSUpgrader initialize a new websocket upgarder that upgrades an incoming http request
// to a websocket connection. It takes in Options that can be used to customize the upgraded connections.
func NewWSUpgrader(opts ...Options) *WSUpgrader {
	defaultUpgrader := &websocket.Upgrader{
		// 默认支持跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	for _, opt := range opts {
		opt(defaultUpgrader)
	}

	return &WSUpgrader{
		Upgrader: defaultUpgrader,
	}
}

func (*Connection) Context() context.Context {
	return context.TODO() // Implement proper context handling if needed
}

func (*Connection) Param(_ string) string {
	return "" // Not applicable for WebSocket, can be implemented if needed
}

func (*Connection) PathParam(_ string) string {
	return "" // Not applicable for WebSocket, can be implemented if needed
}

func (w *Connection) Bind(v interface{}) error {
	_, message, err := w.Conn.ReadMessage()
	if err != nil {
		return err
	}

	switch v := v.(type) {
	case *string:
		*v = string(message)
	default:
		return json.Unmarshal(message, v)
	}

	return nil
}

func (*Connection) HostName() string {
	return "" // Not applicable for WebSocket, can be implemented if needed
}

func (*Connection) Params(string) []string {
	return nil
}

// Dial wraps websocket dialer to establish a new websocket connection.
func (c *Connection) Dial(addr string) (*Connection, error) {
	dialer := websocket.DefaultDialer
	conn, resp, err := dialer.Dial(addr, nil)
	if err != nil {
		return nil, ErrorConnection
	}
	defer conn.Close()
	defer resp.Body.Close()

	return &Connection{Conn: conn}, nil
}

// Manager is a websocket manager that handles the upgrader and manages all
// active connections through ConnectionHub.
type Manager struct {
	ConnectionHub
	WebSocketUpgrader *WSUpgrader
}

// ConnectionHub stores and provide functionality to work with
// all active connections with websocket clients.
type ConnectionHub struct {
	mu                   sync.RWMutex
	WebSocketConnections map[string]*Connection
}

// New intializes a new websocket manager with default websocket upgrader.
func New() *Manager {
	return &Manager{
		WebSocketUpgrader: NewWSUpgrader(),
		ConnectionHub: ConnectionHub{
			mu:                   sync.RWMutex{},
			WebSocketConnections: make(map[string]*Connection),
		},
	}
}

// Upgrade calls the upgrader to upgrade an http connection to a websocket connection.
func (u *WSUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	return u.Upgrader.Upgrade(w, r, responseHeader)
}

// GetWebsocketConnection returns a websocket connection which has been intialized in the middleware.
func (ws *Manager) GetWebsocketConnection(connID string) *Connection {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	return ws.WebSocketConnections[connID]
}

// AddWebsocketConnection add a new connection with the connection id key.
func (ws *Manager) AddWebsocketConnection(connID string, conn *Connection) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	ws.WebSocketConnections[connID] = conn
}

// CloseConnection closes a websocket connection and then removes it from the connection hub.
func (ws *Manager) CloseConnection(connID string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if conn, ok := ws.WebSocketConnections[connID]; ok {
		conn.Close()

		delete(ws.WebSocketConnections, connID)
	}
}