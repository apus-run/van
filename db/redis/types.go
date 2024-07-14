package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

// Option is the database configuration option
type Option func(*Config)

// Config is the database configuration
type Config struct {
	*redis.Options
}

// DefaultOptions .
func DefaultOptions() *Config {
	return &Config{
		&redis.Options{
			Password: "",
			Addr:     "",
			DB:       0,
		},
	}
}

func Apply(opts ...Option) *Config {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}
	return options
}

func WithRedisConfig(f func(options *Config)) Option {
	return func(config *Config) {
		f(config)
	}
}

// UniqKey 用来唯一标识一个Config配置
func (config *Config) UniqKey() string {
	return fmt.Sprintf("%v_%v_%v_%v", config.Addr, config.DB, config.Username, config.Network)
}
