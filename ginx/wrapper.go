package ginx

import (
	"context"
	"net/http"
	"strconv"

	"github.com/apus-run/van/pkg/value"
	"github.com/gin-gonic/gin"
)

// WrapContext returns a context wrapped by this file
func WrapContext(c *gin.Context) *Context {
	return &Context{
		Context: c,
	}
}

// Handle convert HandlerFunc to gin.HandlerFunc
func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := WrapContext(ginCtx)
		h(c)
	}
}

func ProxyHandle(fn ProxyHandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := WrapContext(ginCtx)
		proxy, err := fn(c)

		if err != nil {
			c.Data(http.StatusOK, "text/text", []byte(err.Error()))
			c.Abort()
		} else {
			proxy.ServeHTTP(c.Writer, c.Request)
		}
	}
}

type ginKey struct{}

// NewGinContext returns a new Context that carries gin.Context value.
func NewGinContext(ctx context.Context, c *gin.Context) context.Context {
	return context.WithValue(ctx, ginKey{}, c)
}

// FromGinContext returns the gin.Context value stored in ctx, if any.
func FromGinContext(ctx context.Context) (c *gin.Context, ok bool) {
	c, ok = ctx.Value(ginKey{}).(*gin.Context)
	return
}

func (ctx *Context) Method() string {
	return ctx.Request.Method
}

func (ctx *Context) GetContext() context.Context {
	return ctx.Context.Request.Context()
}

func (ctx *Context) URLPath() string {
	return ctx.FullPath()
}

func (ctx *Context) StatusCode() int {
	return ctx.Writer.Status()
}

func (ctx *Context) BytesWritten() int64 {
	return int64(ctx.Writer.Size())
}

func (c *Context) Param(key string) value.AnyValue {
	return value.AnyValue{
		Value: c.Context.Param(key),
	}
}

func (c *Context) Query(key string) value.AnyValue {
	return value.AnyValue{
		Value: c.Context.Query(key),
	}
}

func (c *Context) Cookie(key string) value.AnyValue {
	val, err := c.Context.Cookie(key)
	return value.AnyValue{
		Value: val,
		Error: err,
	}
}

// JSON returns JSON response
// e.x. {"code":<code>, "msg":<msg>, "data":<data>, "details":<details>}
func (ctx *Context) JSON(httpStatus int, resp Result) {
	ctx.Context.JSON(httpStatus, resp)
}

// JSONOK returns JSON response with successful business code and data
// e.x. {"code": 200, "msg":"成功", "data":<data>}
func (ctx *Context) JSONOK(msg string, data any) {
	j := new(Result)
	j.Code = CodeOK
	j.Msg = msg

	switch d := data.(type) {
	case error:
		j.Data = d.Error()
	case nil:
		j.Data = gin.H{}
	default:
		j.Data = data
	}

	ctx.Context.JSON(http.StatusOK, j)
}

// Success c.Success()
func (ctx *Context) Success(data ...any) {
	j := new(Result)
	j.Code = CodeOK
	j.Msg = "ok"

	if len(data) > 0 {
		j.Data = data[0]
	} else {
		j.Data = ""
	}

	ctx.Context.JSON(http.StatusOK, j)
}

// JSONE returns JSON response with failure business code ,msg and data
// e.x. {"code":<code>, "msg":<msg>, "data":<data>}
// c.JSONE(5, "系统错误", err)
func (ctx *Context) JSONE(code int, msg string, data any) {
	j := new(Result)
	j.Code = code
	j.Msg = msg

	switch d := data.(type) {
	case error:
		j.Data = d.Error()
	case nil:
		j.Data = gin.H{}
	default:
		j.Data = data
	}

	ctx.Context.JSON(http.StatusOK, j)
}

// NotFound 未找到相关路由
func (ctx *Context) NotFound() {
	ctx.String(http.StatusNotFound, "the route not found")
}

// GetClientLocale returns the client locale name
func (ctx *Context) GetClientLocale() string {
	value := ctx.GetHeader(AcceptLanguageHeaderName)

	return value
}

// GetClientTimezoneOffset returns the client timezone offset
func (ctx *Context) GetClientTimezoneOffset() (int16, error) {
	value := ctx.GetHeader(ClientTimezoneOffsetHeaderName)
	offset, err := strconv.Atoi(value)

	if err != nil {
		return 0, err
	}

	return int16(offset), nil
}

// SetRequestId sets the given request id to context
func (ctx *Context) SetRequestId(requestId string) {
	ctx.Set(requestIdFieldKey, requestId)
}

// GetRequestId returns the current request id
func (ctx *Context) GetRequestId() string {
	requestId, exists := ctx.Get(requestIdFieldKey)

	if !exists {
		return ""
	}

	return requestId.(string)
}
