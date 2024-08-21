package conf_test

import (
	"testing"

	"github.com/apus-run/van/conf"
	"github.com/apus-run/van/conf/file"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// 在初始化模块的时候再读配置信息
	type DB struct {
		DSN string `yaml:"dsn"`
	}
	type Config struct {
		DB
	}
	var cfg Config
	var db DB

	t.Run("using yaml config", func(t *testing.T) {
		c := conf.New([]conf.Source{
			file.NewSource("testdata/dev.yaml"),
		})
		err := c.Load()

		assert.NoError(t, err)

		assert.NotNil(t, c)

		err = c.File("dev").UnmarshalKey("db", &db)
		if err != nil {
			t.Fatalf("unmarshal key error: %v", err)
		}
		assert.NoError(t, err)

		t.Logf("db: %+v", db)

		err = c.Scan("dev", &cfg)
		if err != nil {
			t.Fatalf("scan error: %v", err)
		}
		assert.NoError(t, err)

		t.Logf("cfg: %v", cfg)

		cf := c.File("dev")
		err = cf.Unmarshal(&cfg)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		t.Logf("cfg: %v", cfg)
	})
}
