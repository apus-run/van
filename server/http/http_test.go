package http

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	g := gin.Default()
	g.Handle("GET", "/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	server := NewServer(g, WithAddr(":8080"), WithShutdownTimeout(time.Second))

	go func() {
		time.Sleep(2 * time.Second)

		resp, err := http.Get("http://localhost:8080")
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode, http.StatusOK)

		time.Sleep(1 * time.Second)
		err = server.Stop(context.Background())
		assert.NoError(t, err)
		t.Log("shutdown completed")
	}()

	_ = server.Start(context.Background())
}
