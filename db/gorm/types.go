package gorm

import "gorm.io/gorm"

// Driver is the client driver
type Driver int

// The Driver Type of native client
const (
	Unknown Driver = iota
	MySQL
	PostgreSQL
	SQLite
	SQLServer
	ClickHouse
)

// driverMapToString is the safemap of [driver, name]
var driverMapToString = map[Driver]string{
	MySQL:      "mysql",
	PostgreSQL: "postgres",
	SQLite:     "sqlite",
	SQLServer:  "sqlserver",
	ClickHouse: "clickhouse",
	Unknown:    "unknown",
}

var stringMapToDriver = map[string]Driver{
	"mysql":      MySQL,
	"postgres":   PostgreSQL,
	"sqlite":     SQLite,
	"sqlserver":  SQLServer,
	"clickhouse": ClickHouse,
	"unknown":    Unknown,
}

// DriverTypeMap is the safemap of driver [name, driver]
var DriverTypeMap = ReverseMap(driverMapToString)

// String convert the DriverType to string
func (d *Driver) String() string {
	if val, ok := driverMapToString[*d]; ok {
		return val
	}
	return driverMapToString[Unknown]
}

// DriverType convert the string to DriverType
func (d *Driver) DriverType(name string) Driver {
	if val, ok := DriverTypeMap[name]; ok {
		*d = val
		return val
	}
	return Unknown
}

// ReverseMap just reverse the safemap from [key, value] to [value, key]
func ReverseMap[K comparable, V comparable](m map[K]V) map[V]K {
	n := make(map[V]K, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}

// Option 代表初始化的时候的选项
type Option func(*Config)

// Config is the database configuration
type Config struct {
	Driver Driver `json:"driver"`
	DSN    string `json:"dsn"`

	// 以下配置关于gorm
	*gorm.Config // 集成gorm的配置
}

// DefaultOptions .
func DefaultOptions() *Config {
	return &Config{
		Driver: Unknown,
		DSN:    "",
		Config: &gorm.Config{},
	}
}

func Apply(opts ...Option) *Config {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}
	return options
}

func WithDriver(driver Driver) Option {
	return func(config *Config) {
		config.Driver = driver
	}
}

func WithDSN(dsn string) Option {
	return func(config *Config) {
		config.DSN = dsn
	}
}

func WithGormConfig(f func(options *Config)) Option {
	return func(config *Config) {
		f(config)
	}
}

// WithDryRun 设置空跑模式
func WithDryRun() Option {
	return func(config *Config) {
		config.DryRun = true

	}
}

// WithFullSaveAssociations 设置保存时候关联
func WithFullSaveAssociations() Option {
	return func(config *Config) {
		config.FullSaveAssociations = true
	}
}
