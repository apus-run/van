package sqlx

import (
	"context"
	"fmt"
	"sync"

	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/singleflight"
)

var _ Database = (*Helper)(nil)

type Database interface {
	// GetDB 获取数据库连接
	GetDB(ctx context.Context, options ...Option) (*DB, error)

	ConnectDB(ctx context.Context, db *DB) (bool, error)
	CloseDB(ctx context.Context, options ...Option) error
}

type Helper struct {
	lock  *sync.RWMutex
	group *singleflight.Group

	dbs map[string]*DB
}

func NewHelper() *Helper {
	return &Helper{
		lock:  &sync.RWMutex{},
		group: &singleflight.Group{},
		dbs:   make(map[string]*DB),
	}
}

func (h *Helper) GetDB(ctx context.Context, options ...Option) (*DB, error) {
	config := Apply(options...)

	// 判断是否已经实例化了DB
	h.lock.RLock()
	if db, ok := h.dbs[config.DSN]; ok {
		h.lock.RUnlock()
		return db, nil
	}
	h.lock.RUnlock()

	if len(config.DSN) == 0 {
		return nil, fmt.Errorf("database dsn is empty")
	}
	if config.Driver == 0 {
		return nil, fmt.Errorf("unknown database driver for: %q", config.Driver)
	}
	v, err, _ := h.group.Do(config.DSN, func() (any, error) {
		sdb := sqlx.MustOpen(driverMapToString[config.Driver], config.DSN)

		// Mapper function for SQL name mapping, snake_case table names
		sdb.MapperFunc(strcase.ToSnake)

		db := &DB{sdb}

		h.lock.Lock()
		defer h.lock.Unlock()
		h.dbs[config.DSN] = db

		return db, nil
	})

	return v.(*DB), err
}

func (h *Helper) ConnectDB(ctx context.Context, db *DB) (bool, error) {
	if err := db.PingContext(ctx); err != nil {
		return false, err
	}
	return true, nil
}

func (h *Helper) CloseDB(ctx context.Context, options ...Option) error {
	config := Apply(options...)
	if db, ok := h.dbs[config.DSN]; ok {
		if err := db.Close(); err != nil {
			return err
		}
		// 删除数据库实例
		delete(h.dbs, config.DSN)
	}
	return nil
}
