package conf

import "github.com/spf13/viper"

// KV KeyValue is conf key value.
type KV struct {
	Key    string
	Value  []byte
	Format string
	Path   string
}

func (k *KV) Read(p []byte) (n int, err error) {
	return copy(p, k.Value), nil
}

type Source interface {
	Load() ([]*KV, error)
}

type Conf interface {
	File(filename string) *viper.Viper
	Scan(filename string, obj any) error
	Get(filename string, key string) any
	Load() error
	Watch(fn func())
}
