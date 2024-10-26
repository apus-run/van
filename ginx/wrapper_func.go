package ginx

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func W(fn func(ctx *Context) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := fn(&Context{Context: ctx})
		if errors.Is(err, ErrNoResponse) {
			slog.Debug("不需要响应", slog.Any("err", err))
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			slog.Debug("未授权", slog.Any("err", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err != nil {
			slog.Error("执行业务逻辑失败", slog.Any("err", err))
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func B[Req any](fn func(ctx *Context, req Req) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			slog.Debug("绑定参数失败", slog.Any("err", err))
			return
		}
		res, err := fn(&Context{Context: ctx}, req)
		if errors.Is(err, ErrNoResponse) {
			slog.Debug("不需要响应", slog.Any("err", err))
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			slog.Debug("未授权", slog.Any("err", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err != nil {
			slog.Error("执行业务逻辑失败", slog.Any("err", err))
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WC(fn func(*gin.Context, func() jwt.Claims) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawVal, ok := ctx.Get("claims")

		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			slog.Error("无法获得 claims",
				slog.String("path", ctx.Request.URL.Path))
			return
		}

		claims, ok := rawVal.(func() jwt.Claims)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			slog.Error("无法获得 claims",
				slog.String("path", ctx.Request.URL.Path))
			return
		}

		res, err := fn(ctx, claims)
		if err != nil {
			slog.Error("执行业务逻辑失败",
				slog.Any("err", err))
		}

		// TODO 可以在这里放一些可观测性的中间件

		ctx.JSON(http.StatusOK, res)
	}
}

func BC[Req any](fn func(*gin.Context, Req, func() jwt.Claims) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			slog.Error("解析请求失败", slog.Any("err", err))
			return
		}

		rawVal, ok := ctx.Get("claims")

		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			slog.Error("无法获得 claims",
				slog.String("path", ctx.Request.URL.Path))
			return
		}

		claims, ok := rawVal.(func() jwt.Claims)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			slog.Error("无法获得 claims",
				slog.String("path", ctx.Request.URL.Path))
			return
		}

		res, err := fn(ctx, req, claims)

		// TODO 可以在这里放一些可观测性的中间件

		if err != nil {
			slog.Error("执行业务逻辑失败", slog.Any("err", err))
		}

		ctx.JSON(http.StatusOK, res)
	}
}
