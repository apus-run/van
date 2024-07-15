package redis

import (
	"context"
	"testing"
	"time"
)

func TestRedis_GetClient(t *testing.T) {
	ctx := context.Background()
	h := NewHelper()
	client, err := h.GetDB(context.Background(), WithRedisConfig(func(options *Config) {
		options.Addr = "localhost:6379"
		options.DB = 0
		options.Username = "root"
	}))
	if err != nil {
		t.Fatal(err)
	}

	// 检测数据库是否可以连接
	cmd := client.Ping(ctx)
	if cmd.Err() != nil {
		t.Fatal(cmd.Err())
	}

	err = client.Set(ctx, "foo", "bar", 1*time.Hour).Err()

	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("foo", val)

	err = client.Del(ctx, "foo").Err()
	if err != nil {
		t.Fatal(err)
	}

}
