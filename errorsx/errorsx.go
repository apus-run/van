package errorsx

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	httpstatus "github.com/apus-run/van/server/http/status"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

var (
	// EnableStackCapture 堆栈捕获开关（生产环境可关闭）
	EnableStackCapture = true
	// StackDepth 堆栈捕获深度
	StackDepth = 5
)

// Error 定义了项目体系中使用的错误类型，用于描述错误的详细信息.
type Error struct {
	// 基本信息
	Code    int    `json:"code"`    // HTTP 状态码
	Reason  string `json:"reason"`  // 业务错误码
	Message string `json:"message"` // 给用户看的错误信息

	// 额外信息
	Metadata map[string]string `json:"metadata,omitempty"` // 附加的元数据，通常用于提供额外的上下文或调试信息
	Cause    error             `json:"cause,omitempty"`    // 原始错误信息，通常用于记录日志或调试
	Stack    []string          `json:"stack,omitempty"`    // 错误发生时的调用栈信息，通常用于调试和排查问题
}

func New(code int, reason string) *Error {
	return &Error{
		Code:   code,
		Reason: reason,
	}
}

// Error 实现 error 接口中的 `Error` 方法.
func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d, reason = %s, message = %s", e.Code, e.Reason, e.Message)
}

// WithMetadata 用于设置与错误相关的元信息，通常用于提供额外的上下文或调试信息.
func (e *Error) WithMetadata(md map[string]string) *Error {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}

	// 遍历元数据并添加到错误对象中
	// 如果元数据中已经存在相同的键，则保留原有的值
	for k, v := range md {
		e.Metadata[k] = v
	}
	return e
}

// KV 使用 key-value 对设置元数据.
func (e *Error) KV(kvs ...string) *Error {
	if len(kvs)%2 != 0 {
		return e // 忽略不完整键值对
	}

	// 如果元数据为空，则初始化
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}

	// 遍历键值对并添加到元数据中
	// 如果元数据中已经存在相同的键，则保留原有的值
	for i := 0; i < len(kvs); i += 2 {
		// kvs 必须是成对的
		if i+1 < len(kvs) {
			e.Metadata[kvs[i]] = kvs[i+1]
		}
	}
	return e
}

// WithMessage 用于设置错误的详细信息，通常用于提供更具体的错误描述.
func (e *Error) WithMessage(msg string) *Error {
	e.Message = msg
	return e
}

// WithCause with original error
func (e *Error) WithCause(err error) *Error {
	e.Cause = err
	return e
}

// WithStack with stack
func (e *Error) WithStack() *Error {
	// 如果启用了堆栈捕获，则捕获当前堆栈信息
	// skip: 2 表示跳过当前函数和 runtime.Callers
	// depth: StackDepth 表示捕获的堆栈深度
	// 如果 StackDepth 为 0，则表示无限制捕获
	if EnableStackCapture {
		e.Stack = captureStack(2, StackDepth)
	}
	return e
}

// WithRequestID 设置请求 ID.
func (e *Error) WithRequestID(requestID string) *Error {
	return e.KV("X-Request-ID", requestID) // 设置请求 ID
}

// WithUserID 设置用户 ID.
func (e *Error) WithUserID(userID string) *Error {
	return e.KV("X-User-ID", userID) // 设置用户 ID
}

func (e *Error) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		str := bytes.NewBuffer([]byte{})
		str.WriteString(fmt.Sprintf("code: %d, ", e.Code))
		str.WriteString("reason: ")
		str.WriteString(e.Reason + ", ")
		str.WriteString("message: ")
		str.WriteString(e.Message)
		if len(e.Metadata) > 0 {
			str.WriteString(", metadata: ")
			for k, v := range e.Metadata {
				str.WriteString(fmt.Sprintf("%s=%s ", k, v))
			}
		}
		if e.Cause != nil {
			str.WriteString(", error: ")
			str.WriteString(e.Cause.Error())
		}
		if len(e.Stack) > 0 {
			str.WriteString(", stack: ")
			for _, s := range e.Stack {
				str.WriteString(fmt.Sprintf("%s\n", s))
			}
		}

		fmt.Fprintf(state, "%s", strings.Trim(str.String(), "\r\n\t"))
	default:
		fmt.Fprintf(state, e.Message)
	}
}

// GRPCStatus 返回 gRPC 状态表示.
func (e *Error) GRPCStatus() *status.Status {
	st := status.New(
		httpstatus.ToGRPCCode(e.Code),
		fmt.Sprintf("%s: %s", e.Reason, e.Message),
	)

	// 添加错误详情
	details := &errdetails.ErrorInfo{
		Reason:   e.Reason,
		Metadata: e.Metadata,
	}

	st, _ = st.WithDetails(details)

	return st
}

// Unwrap 返回原始错误.
func (e *Error) Unwrap() error {
	if e.Cause != nil {
		return e.Cause
	}
	return nil
}

// Is 判断当前错误是否与目标错误匹配.
// 它会递归遍历错误链，并比较 Error 实例的 Code 和 Reason 字段.
// 如果 Code 和 Reason 均相等，则返回 true；否则返回 false.
func (e *Error) Is(target error) bool {
	if targetErr := new(Error); errors.As(target, &targetErr) {
		return targetErr.Code == e.Code && targetErr.Reason == e.Reason
	}
	return errors.Is(e.Cause, target)
}

func (e *Error) As(target any) bool {
	if t, ok := target.(**Error); ok {
		*t = e
		return true
	}
	return false
}

func (e *Error) Clone() *Error {
	return &Error{
		Code:     e.Code,
		Reason:   e.Reason,
		Message:  e.Message,
		Metadata: e.Metadata,
		Cause:    e.Cause,
		Stack:    e.Stack,
	}
}

// Code 返回错误的 HTTP 代码.
func Code(err error) int {
	if err == nil {
		return http.StatusOK //nolint:mnd
	}
	return FromError(err).Code
}

// Reason 返回特定错误的原因.
func Reason(err error) string {
	if err == nil {
		return ""
	}
	return FromError(err).Reason
}

// FromError 尝试将一个通用的 error 转换为自定义的 *Error 类型.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}

	// 处理自定义错误类型
	var target *Error
	if errors.As(err, &target) {
		return target
	}

	// 处理 gRPC 错误
	if st, ok := status.FromError(err); ok {
		return fromGRPCStatus(st)
	}

	// 处理标准错误
	return &Error{
		Code:    http.StatusInternalServerError,
		Reason:  "InternalError",
		Message: err.Error(),
		Cause:   err,
		Stack:   captureStack(2, 5),
	}
}

// 从 gRPC 状态转换
func fromGRPCStatus(st *status.Status) *Error {
	code := httpstatus.FromGRPCCode(st.Code())
	e := New(code, "GRPCError").WithMessage(st.Message())

	// 提取详情信息
	for _, detail := range st.Details() {
		if info, ok := detail.(*errdetails.ErrorInfo); ok {
			e.Reason = info.Reason
			e.Metadata = info.Metadata
		}
	}

	return e
}

// captureStack 增强的堆栈捕获方法
// skip: 跳过的调用层级（通常为 2：跳过自身和 runtime.Callers）
// depth: 最大捕获深度（0 表示无限制）
func captureStack(skip, depth int) []string {
	pcs := make([]uintptr, depth)
	n := runtime.Callers(skip+1, pcs)
	if n == 0 {
		return nil
	}

	frames := runtime.CallersFrames(pcs[:n])
	stack := make([]string, 0, n)
	for {
		frame, more := frames.Next()
		stack = append(stack,
			fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}
	return stack
}
