package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/apus-run/van/store/logger/empty"
	"github.com/apus-run/van/store/where"
)

// DBProvider defines an interface for providing a database connection.
type DBProvider interface {
	// DB returns the database instance for the given context.
	DB(ctx context.Context, wheres ...where.Where) *gorm.DB
}

// Option defines a function type for configuring the Store.
type Option[T any] func(*Store[T])

// Store represents a generic data store with logging capabilities.
type Store[T any] struct {
	logger  Logger
	storage DBProvider
}

// WithLogger returns an Option function that sets the provided Logger to the Store for logging purposes.
func WithLogger[T any](logger Logger) Option[T] {
	return func(s *Store[T]) {
		s.logger = logger
	}
}

// NewStore creates a new instance of Store with the provided DBProvider.
func NewStore[T any](storage DBProvider, logger Logger) *Store[T] {
	if logger == nil {
		logger = empty.NewLogger()
	}

	return &Store[T]{
		logger:  logger,
		storage: storage,
	}
}

// db retrieves the database instance and applies the provided where conditions.
func (s *Store[T]) db(ctx context.Context, wheres ...where.Where) *gorm.DB {
	dbInstance := s.storage.DB(ctx)
	for _, whr := range wheres {
		if whr != nil {
			dbInstance = whr.Where(dbInstance)
		}
	}
	return dbInstance
}

// Create inserts a new object into the database.
func (s *Store[T]) Create(ctx context.Context, obj *T) error {
	if err := s.db(ctx).Create(obj).Error; err != nil {
		s.logger.Error(ctx, err, "Failed to insert object into database", "object", obj)
		return err
	}
	return nil
}

// Update modifies an existing object in the database.
func (s *Store[T]) Update(ctx context.Context, obj *T) error {
	if err := s.db(ctx).Save(obj).Error; err != nil {
		s.logger.Error(ctx, err, "Failed to update object in database", "object", obj)
		return err
	}
	return nil
}

// Delete removes an object from the database based on the provided where options.
func (s *Store[T]) Delete(ctx context.Context, opts *where.Options) error {
	err := s.db(ctx, opts).Delete(new(T)).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error(ctx, err, "Failed to delete object from database", "conditions", opts)
		return err
	}
	return nil
}

// Get retrieves a single object from the database based on the provided where options.
func (s *Store[T]) Get(ctx context.Context, opts *where.Options) (*T, error) {
	var obj T
	if err := s.db(ctx, opts).First(&obj).Error; err != nil {
		s.logger.Error(ctx, err, "Failed to retrieve object from database", "conditions", opts)
		return nil, err
	}
	return &obj, nil
}

// List retrieves a list of objects from the database based on the provided where options.
func (s *Store[T]) List(ctx context.Context, opts *where.Options) (count int64, ret []*T, err error) {
	err = s.db(ctx, opts).Order("id desc").Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		s.logger.Error(ctx, err, "Failed to list objects from database", "conditions", opts)
	}
	return
}

// Count returns the number of objects in the database that match the provided where options.
func (s *Store[T]) Count(ctx context.Context, opts *where.Options) (int64, error) {
	var count int64
	err := s.db(ctx, opts).Model(new(T)).Count(&count).Error
	if err != nil {
		s.logger.Error(ctx, err, "Failed to count objects in database", "conditions", opts)
		return 0, err
	}
	return count, nil
}

// Transaction executes a function within a database transaction.
func (s *Store[T]) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := s.db(ctx).Begin()
	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			s.logger.Error(ctx, err, "Transaction rollback failed", "error", rollbackErr)
		}
		s.logger.Error(ctx, err, "Transaction failed", "error", err)
		return err
	}
	return tx.Commit().Error
}
