package gorm

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/singleflight"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var _ Database = (*Helper)(nil)

type Database interface {
	GetDB(ctx context.Context, options ...Option) (*gorm.DB, error)
	ConnectDB(ctx context.Context, db *gorm.DB) (bool, error)
	CloseDB(ctx context.Context, options ...Option) error
}

// Transaction 事物接口
type Transaction interface {
	//Execute 执行一个事务方法
	//该方法接受一个上下文对象 ctx 和一个函数 fn，该函数用于执行事务内的操作。
	//方法返回一个错误，可能是事务执行过程中出现的任何错误。
	Execute(ctx context.Context, fn func(ctx context.Context) error) error
}

type DB = gorm.DB

type Helper struct {
	lock  *sync.RWMutex
	group *singleflight.Group

	dbs map[string]*gorm.DB
}

func NewHelper() *Helper {
	return &Helper{
		lock:  &sync.RWMutex{},
		group: &singleflight.Group{},
		dbs:   make(map[string]*gorm.DB),
	}
}

func (h *Helper) GetDB(ctx context.Context, options ...Option) (*gorm.DB, error) {
	config := Apply(options...)

	// 判断是否已经实例化了gorm.DB
	h.lock.RLock()
	if db, ok := h.dbs[config.DSN]; ok {
		h.lock.RUnlock()
		return db, nil
	}
	h.lock.RUnlock()

	v, err, _ := h.group.Do(config.DSN, func() (any, error) {
		var db *gorm.DB
		var err error

		if len(config.DSN) == 0 {
			return nil, errors.New("database dsn is empty")
		}
		switch config.Driver {
		case MySQL:
			db, err = gorm.Open(mysql.Open(config.DSN), config)
		case PostgreSQL:
			db, err = gorm.Open(postgres.Open(config.DSN), config.Config)
		case SQLite:
			db, err = gorm.Open(sqlite.Open(config.DSN), config.Config)
		case SQLServer:
			db, err = gorm.Open(sqlserver.Open(config.DSN), config.Config)
		case ClickHouse:
			db, err = gorm.Open(clickhouse.Open(config.DSN), config.Config)
		default:
			return nil, errors.New("unknown database driver")
		}

		if err != nil {
			return nil, fmt.Errorf("open database error: %w", err)
		}

		h.lock.Lock()
		defer h.lock.Unlock()
		h.dbs[config.DSN] = db

		return db, nil
	})

	if err != nil {
		return nil, err
	}
	return v.(*gorm.DB), err
}

func (h *Helper) ConnectDB(ctx context.Context, db *gorm.DB) (bool, error) {
	sqlDb, err := db.DB()
	if err != nil {
		return false, fmt.Errorf("CanConnect Ping error: %w", err)
	}
	if err := sqlDb.Ping(); err != nil {
		return false, fmt.Errorf("CanConnect Ping error: %w", err)
	}
	return true, nil
}

func (h *Helper) CloseDB(ctx context.Context, options ...Option) error {
	config := Apply(options...)
	if db, ok := h.dbs[config.DSN]; ok {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		if err := sqlDB.Close(); err != nil {
			return err
		}

		// 删除数据库实例
		delete(h.dbs, config.DSN)
	}
	return nil
}
