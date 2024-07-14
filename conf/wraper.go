package conf

import (
	"errors"
	"log"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	files  []Source
	cached *sync.Map
}

func New(files []Source) *Config {
	conf := &Config{
		files:  files,
		cached: &sync.Map{},
	}

	err := conf.Load()
	if err != nil {
		return nil
	}

	return conf
}

func (c *Config) Watch(fn func()) {
	c.cached.Range(func(key, value any) bool {
		v := value.(*viper.Viper)
		v.OnConfigChange(func(e fsnotify.Event) {
			fn()
		})
		v.WatchConfig()
		return true
	})
}

func (c *Config) Scan(filename string, obj any) error {
	err := c.File(filename).Unmarshal(obj)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) File(filename string) *viper.Viper {
	if v, ok := c.cached.Load(filename); ok {
		return v.(*viper.Viper)
	}
	return nil
}

func (c *Config) Get(filename string, key string) any {
	return c.File(filename).Get(key)
}

func (c *Config) Load() error {
	if len(c.files) == 0 {
		return nil
	}
	for _, file := range c.files {
		kvs, err := file.Load()
		if err != nil {
			return err
		}

		for _, kv := range kvs {
			v := viper.New()
			v.SetConfigType(kv.Format)
			v.SetConfigFile(kv.Path)

			if err := v.ReadInConfig(); err != nil {
				var configFileNotFoundError viper.ConfigFileNotFoundError
				if errors.As(err, &configFileNotFoundError) {
					log.Printf("Using conf file: %s [%s]\n", viper.ConfigFileUsed(), err)
					return errors.New("conf file not found")
				}
				return err
			}
			v.AutomaticEnv()

			name := strings.TrimSuffix(path.Base(kv.Key), filepath.Ext(kv.Key))
			c.cached.Store(name, v)
		}
	}
	return nil
}
