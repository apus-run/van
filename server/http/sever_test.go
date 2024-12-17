package http_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/apus-run/van/server"
	httpServer "github.com/apus-run/van/server/http"
)

func TestNewServer(t *testing.T) {
	g := gin.Default()
	g.Handle("GET", "/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	srv := httpServer.NewServer(server.WithHandler(g), server.WithAddress(":8080"))

	// wsManager := ws.New()
	// wsUpgrader := ws.NewWSUpgrader(
	// 	ws.WithHandshakeTimeout(5*time.Second), // Set handshake timeout
	// 	ws.WithReadBufferSize(2048),            // Set read buffer size
	// 	ws.WithWriteBufferSize(2048),           // Set write buffer size
	// 	ws.WithSubprotocols("chat", "binary"),  // Specify subprotocols
	// 	ws.WithCompression(),                   // Enable compression
	// )

	// server.ws = wsManager
	// server.ws.WebSocketUpgrader.Upgrader = wsUpgrader

	go func() {
		time.Sleep(2 * time.Second)

		resp, err := http.Get("http://localhost:8080")
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode, http.StatusOK)

		time.Sleep(1 * time.Second)
		err = srv.Stop(context.Background())
		assert.NoError(t, err)
		t.Log("shutdown completed")
	}()

	_ = srv.Start(context.Background())
}

func TestServer(t *testing.T) {
	ctx := context.Background()
	g := gin.Default()
	g.Handle("GET", "/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	srv := httpServer.NewServer(server.WithHandler(g), server.WithAddress(":8080"))

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(srv, done)

	// Start the server
	err := srv.Start(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the shutdown to be complete
	<-done

	t.Logf("Graceful shutdown complete.")
}

func gracefulShutdown(srv server.Server, done chan bool) {
	// graceful shutdown
	exitSignals := []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	}

	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), exitSignals...)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func TestServerWG(t *testing.T) {
	ctx := context.Background()
	wg := sync.WaitGroup{}

	g := gin.Default()
	g.Handle("GET", "/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	srv := httpServer.NewServer(server.WithHandler(g), server.WithAddress(":8080"))

	// graceful shutdown
	exitSignals := []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	}

	// 等待中断信号来优雅地关闭服务器
	quit := make(chan os.Signal, len(exitSignals))
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, exitSignals...)

	wg.Add(1)
	go func() {
		<-quit
		stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Stop(stopCtx); err != nil {
			log.Printf("Server forced to shutdown with error: %v", err)
		}
		wg.Done()
	}()

	// Start the server
	err := srv.Start(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	wg.Wait()

	t.Logf("Graceful shutdown complete.")
}
