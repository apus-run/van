package redis

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

var _ Database = (*Helper)(nil)

type (
	// Options is an alias of redis.Options.
	Options = redis.Options
	// Client is an alias of redis.Client.
	Client = redis.Client
	// UniversalClient is an alias of redis.UniversalClient.
	UniversalClient = redis.UniversalClient
	// Tx is an alias of redis.Tx.
	Tx = redis.Tx
	// Cmdable is an alias of redis.Cmdable.
	Cmdable = redis.Cmdable
	// Pipeline is an alias of redis.Pipeline.
	Pipeline = redis.Pipeline

	// StatusCmd is an alias of redis.StatusCmd.
	StatusCmd = redis.StatusCmd
	// StringSliceCmd is an alias of redis.StringSliceCmd.
	StringSliceCmd = redis.StringSliceCmd
	// BoolCmd is an alias of redis.BoolCmd.
	BoolCmd = redis.BoolCmd
	// IntCmd is an alias of redis.IntCmd.
	IntCmd = redis.IntCmd
	// FloatCmd is an alias of redis.FloatCmd.
	FloatCmd = redis.FloatCmd
	// StringCmd is an alias of redis.StringCmd.
	StringCmd = redis.StringCmd
	// Script is an alias of redis.Script.
	Script = redis.Script
)

type Database interface {
	// GetDB 获取数据库连接
	GetDB(ctx context.Context, options ...Option) (redis.Cmdable, error)

	ConnectDB(ctx context.Context, db redis.Cmdable) (bool, error)
	CloseDB(ctx context.Context, options ...Option) error
}

// Transaction 事物接口
type Transaction interface {
	// Execute 执行一个事务方法，func为一个需要保证事务完整性的业务方法
	Execute(ctx context.Context, fn func(ctx context.Context) error) error
}

type Helper struct {
	lock  *sync.RWMutex
	group *singleflight.Group

	clients map[string]redis.Cmdable
}

func NewHelper() *Helper {
	return &Helper{
		lock:    &sync.RWMutex{},
		group:   &singleflight.Group{},
		clients: make(map[string]redis.Cmdable),
	}
}

func (h *Helper) GetDB(ctx context.Context, options ...Option) (redis.Cmdable, error) {
	config := Apply(options...)

	// 如果最终的config没有设置dsn,就生成dsn
	key := config.UniqKey()

	// 判断是否已经实例化了 redis.Client
	h.lock.RLock()
	if db, ok := h.clients[key]; ok {
		h.lock.RUnlock()
		return db, nil
	}
	h.lock.RUnlock()

	v, err, _ := h.group.Do(key, func() (any, error) {
		// 实例化redis.NewClient
		client := redis.NewClient(config.Options)

		h.lock.Lock()
		defer h.lock.Unlock()
		// 挂载到map中，结束配置
		h.clients[key] = client

		return client, nil
	})

	return v.(redis.Cmdable), err
}

func (h *Helper) ConnectDB(ctx context.Context, db redis.Cmdable) (bool, error) {
	err := db.Ping(ctx).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *Helper) CloseDB(ctx context.Context, options ...Option) error {
	config := Apply(options...)
	if rdb, ok := h.clients[config.UniqKey()]; ok {
		err := rdb.Shutdown(ctx).Err()
		if err != nil {
			return err
		}
		delete(h.clients, config.UniqKey())
	}
	return nil
}

// NewScript returns a new Script instance.
func NewScript(script string) *Script {
	return redis.NewScript(script)
}
