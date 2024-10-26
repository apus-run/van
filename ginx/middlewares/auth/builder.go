package auth

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/apus-run/van/ginx"
	// jwtx "github.com/apus-run/van/anthx/jwt"
)

// Builder 鉴权，验证用户token是否有效
type Builder struct {
	// 白名单路由地址集合, 放行
	whitePathList []string
}

func NewBuilder() *Builder {
	return &Builder{
		whitePathList: []string{},
	}
}

func (b *Builder) IgnorePaths(whitePath string) *Builder {
	b.whitePathList = append(b.whitePathList, whitePath)
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return ginx.Handle(func(ctx *ginx.Context) {
		// 白名单路由放行
		for _, path := range b.whitePathList {
			if strings.Contains(ctx.Request.URL.Path, path) {
				ctx.Next()
				return
			}
		}

		// tokenString, err := getJwtFromHeader(ctx)
		// if err != nil {
		// 	ctx.JSONE(http.StatusUnauthorized, "invalid token", nil)
		// 	ctx.Abort()
		// 	return
		// }

		ctx.Next()
	})
}

func getJwtFromHeader(ctx *ginx.Context) (string, error) {
	// 读取请求头的 token
	tokenString := ctx.GetHeader("Authorization")
	if len(tokenString) == 0 {
		return "", errors.New("token 为空")
	}
	strs := strings.SplitN(tokenString, " ", 2)
	if len(strs) != 2 || strs[0] != "Bearer" {
		return "", errors.New("token 不符合规则, Bearer 开头")
	}
	return strs[1], nil
}
