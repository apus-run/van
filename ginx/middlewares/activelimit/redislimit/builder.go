package redislimit

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/atomic"
)

type RedisActiveLimit struct {
	//最大限制个数
	maxActive *atomic.Int64

	//用来记录连接数的key
	key   string
	cmd   redis.Cmdable
	logFn func(msg any, args ...any)
}

// NewRedisActiveLimit 全局限流
func NewRedisActiveLimit(cmd redis.Cmdable, maxActive int64, key string) *RedisActiveLimit {
	return &RedisActiveLimit{
		maxActive: atomic.NewInt64(maxActive),
		key:       key,
		cmd:       cmd,
		logFn: func(msg any, args ...any) {
			fmt.Println(fmt.Sprintf("%v  详细信息: %v", msg, args))
		},
	}
}

func (a *RedisActiveLimit) SetMaxActive(maxActive int64) *RedisActiveLimit {
	a.maxActive.Store(maxActive)
	return a
}

func (a *RedisActiveLimit) SetLogFunc(fun func(msg any, args ...any)) *RedisActiveLimit {
	a.logFn = fun
	return a
}

func (a *RedisActiveLimit) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentCount, err := a.cmd.Incr(ctx, a.key).Result()
		if err != nil {
			//为了安全性 直接返回异常
			a.logFn("redis 加一操作", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer func() {
			if err = a.cmd.Decr(ctx, a.key).Err(); err != nil {
				a.logFn("redis 减一操作", err)
			}
		}()
		if currentCount <= a.maxActive.Load() {
			ctx.Next()
		} else {
			a.logFn("web server ", "限流中..")
			ctx.AbortWithStatus(http.StatusTooManyRequests)
		}
		return
	}
}
