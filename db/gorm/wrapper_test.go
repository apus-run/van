package gorm

import (
	"context"
	"fmt"
	"testing"
)

type myUser struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"type:varchar(200)"`
}

func (myUser) TableName() string {
	return "User"
}

func TestDatabase_OpenDB(t *testing.T) {
	ctx := context.Background()
	h := NewHelper()

	db, err := h.GetDB(ctx, WithGormConfig(func(options *Config) {
		options.Driver = MySQL
		options.DSN = "root:123456@tcp(localhost:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
		options.DisableForeignKeyConstraintWhenMigrating = true

		tag := options.Driver.String()
		if tag == "mysql" {
			t.Log("mysql")
		}
	}))

	if err != nil {
		t.Fatal(err)
	}

	// 检测数据库是否可以连接
	ok, err := h.ConnectDB(ctx, db)
	if err != nil || !ok {
		fmt.Println("数据库连接失败，请检查配置")
		t.Fatal(err)
	}
}
