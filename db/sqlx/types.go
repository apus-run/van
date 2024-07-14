package sqlx

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

// Option is the database configuration option
type Option func(*Config)

// Config is the database configuration
type Config struct {
	Driver Driver `json:"driver"`
	DSN    string `json:"dsn"`
}

// DefaultOptions .
func DefaultOptions() *Config {
	return &Config{
		Driver: Unknown,
		DSN:    "",
	}
}

func Apply(opts ...Option) *Config {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}
	return options
}
