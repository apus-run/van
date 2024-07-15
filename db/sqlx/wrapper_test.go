package sqlx

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestDatabase_OpenDB(t *testing.T) {
	ctx := context.Background()
	h := NewHelper()

	db, err := h.GetDB(ctx, WithDriver(MySQL), WithDSN("root:123456@tcp(localhost:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"))

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
