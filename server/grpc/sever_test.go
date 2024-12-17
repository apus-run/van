package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/apus-run/van/server"
)

func TestNewServer(t *testing.T) {
	addr := "0.0.0.0:9090"
	server := NewServer(server.WithAddress(addr))

	go func() {
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
		defer cancelFunc()

		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		assert.NoError(t, err)

		time.Sleep(1 * time.Second)
		err = server.Stop(ctx)
		assert.NoError(t, err)
		t.Log("shutdown completed")
	}()

	_ = server.Start(context.Background())
}
