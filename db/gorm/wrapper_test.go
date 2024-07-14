package gorm

import (
	"testing"

	"github.com/xo/dburl"
)

func TestDbUrl(t *testing.T) {
	u, err := dburl.Parse("mysql://root:123456@localhost:3306/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", u.DSN)
	t.Logf("%v", u.Driver)
	t.Logf("%v", u.UnaliasedDriver)
}
