package errorsx

import "net/http"

// Success new Success error
func Success(reason string) *Error {
	return New(http.StatusOK, reason).WithMessage("OK")
}
func IsSuccess(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusOK &&
		err.Message == "OK"
}

// BadRequest new BadRequest error
func BadRequest(reason string) *Error {
	return New(http.StatusBadRequest, reason).WithMessage("Bad Request")
}

// IsBadRequest determines if err is BadRequest error.
func IsBadRequest(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusBadRequest &&
		err.Message == "Bad Request"
}

// Unauthorized new Unauthorized error
func Unauthorized(reason string) *Error {
	return New(http.StatusUnauthorized, reason).WithMessage("Unauthorized")
}

// IsUnauthorized determines if err is Unauthorized error.
func IsUnauthorized(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusUnauthorized &&
		err.Message == "Unauthorized"
}

// Forbidden new Forbidden error
// 未授权（已认证但无权限）
func Forbidden(reason string) *Error {
	return New(http.StatusForbidden, reason).WithMessage("Forbidden")
}

// IsForbidden determines if err is Forbidden error.
func IsForbidden(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusForbidden &&
		err.Message == "Forbidden"
}

// NotFound new NotFound error
func NotFound(reason string) *Error {
	return New(http.StatusNotFound, reason).WithMessage("Not Found")
}

// IsNotFound determines if err is NotFound error.
func IsNotFound(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusNotFound &&
		err.Message == "Not Found"
}

// Conflict new Conflict error
func Conflict(reason string) *Error {
	return New(http.StatusConflict, reason).WithMessage("Conflict")
}

// IsConflict determines if err is Conflict error.
func IsConflict(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusConflict &&
		err.Message == "Conflict"
}

// InternalServer new InternalServer error
func InternalServer(reason string) *Error {
	return New(http.StatusInternalServerError, reason).WithMessage("Internal Server Error")
}

// IsInternalServer determines if err is InternalServer error.
func IsInternalServer(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusInternalServerError &&
		err.Message == "Internal Server Error"
}

// ServiceUnavailable new ServiceUnavailable error
func ServiceUnavailable(reason string) *Error {
	return New(http.StatusServiceUnavailable, reason).WithMessage("Service Unavailable")
}

// IsServiceUnavailable determines if err is ServiceUnavailable error.
func IsServiceUnavailable(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusServiceUnavailable &&
		err.Message == "Service Unavailable"
}

// GatewayTimeout new GatewayTimeout error
func GatewayTimeout(reason string) *Error {
	return New(http.StatusGatewayTimeout, reason).WithMessage("Gateway Timeout")
}

// IsGatewayTimeout determines if err is GatewayTimeout error.
func IsGatewayTimeout(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusGatewayTimeout &&
		err.Message == "Gateway Timeout"
}

// ClientClosed new ClientClosed error
// 注意：499 是Nginx定义的非标准状态码（客户端主动关闭连接）
func ClientClosed(reason string) *Error {
	return New(499, reason).WithMessage("Client Closed")
}

// IsClientClosed determines if err is ClientClosed error.
func IsClientClosed(err *Error) bool {
	return err != nil &&
		err.Code == 499 &&
		err.Message == "Client Closed"
}

func TooManyRequests(reason string) *Error {
	return New(http.StatusTooManyRequests, reason).WithMessage("Too Many Requests")
}

func IsTooManyRequests(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusTooManyRequests &&
		err.Message == "Too Many Requests"
}

// TokenInvalid token 无效
func TokenInvalid(reason string) *Error {
	return New(http.StatusUnauthorized, reason).WithMessage("Token Invalid")
}

func IsTokenInvalid(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusUnauthorized &&
		err.Message == "Token Invalid"
}

// TokenExpired token 过期
func TokenExpired(reason string) *Error {
	return New(http.StatusUnauthorized, reason).WithMessage("Token Expired")
}

func IsTokenExpired(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusUnauthorized &&
		err.Message == "Token Expired"
}

// TokenInvalidSignature token 签名无效
// 通常用于处理 JWT 或其他类型的令牌验证失败的情况
func TokenInvalidSignature(reason string) *Error {
	return New(http.StatusUnauthorized, reason).WithMessage("Token Invalid Signature")
}

func IsTokenInvalidSignature(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusUnauthorized &&
		err.Message == "Token Invalid Signature"
}

// Bind bind error
// 绑定参数错误，通常用于处理请求体或查询参数的解析错误
func BindError(reason string) *Error {
	return New(http.StatusBadRequest, reason).WithMessage("Bind Error")
}

func IsBindError(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusBadRequest &&
		err.Message == "Bind Error"
}

// InvalidArguments invalid arguments
// 无效参数错误（函数/方法调用时的参数错误）
func InvalidArguments(reason string) *Error {
	return New(http.StatusBadRequest, reason).WithMessage("Invalid Arguments")
}

func IsInvalidArguments(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusBadRequest &&
		err.Message == "Invalid Arguments"
}

// InvalidParams invalid params
// 无效参数错误（外部输入参数验证失败）
func InvalidParams(reason string) *Error {
	return New(http.StatusBadRequest, reason).WithMessage("Invalid Params")
}

func IsInvalidParams(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusBadRequest &&
		err.Message == "Invalid Params"
}

// Panic panic error
// 通常用于处理程序运行时的异常情况（如空指针、数组越界）
func PanicError(reason string) *Error {
	return New(http.StatusInternalServerError, reason).WithMessage("Panic Error")
}

func IsPanicError(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusInternalServerError &&
		err.Message == "Panic Error"
}

// PageNotFound page not found
// 通常用于处理请求的页面或资源未找到的情况
func PageNotFound(reason string) *Error {
	return New(http.StatusNotFound, reason).WithMessage("Page Not Found")
}

func IsPageNotFound(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusNotFound &&
		err.Message == "Page Not Found"
}

// DBReadError db read error
// 包含 SELECT/FETCH 等查询操作失败
func DBReadError(reason string) *Error {
	return New(http.StatusInternalServerError, reason).WithMessage("DB Read Error")
}

func IsDBReadError(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusInternalServerError &&
		err.Message == "DB Read Error"
}

// DBWriteError db write error
// 包含 INSERT/UPDATE/DELETE 等写入操作失败
func DBWriteError(reason string) *Error {
	return New(http.StatusInternalServerError, reason).WithMessage("DB Write Error")
}

func IsDBWriteError(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusInternalServerError &&
		err.Message == "DB Write Error"
}

// DBTransactionError db transaction error
// 数据库事务操作（BEGIN/COMMIT/ROLLBACK）失败
func DBTransactionError(reason string) *Error {
	return New(http.StatusInternalServerError, reason).WithMessage("DB Transaction Error")
}

func IsDBTransactionError(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusInternalServerError &&
		err.Message == "DB Transaction Error"
}

// PermissionDenied 表示请求没有权限
func PermissionDenied(reason string) *Error {
	return New(http.StatusForbidden, reason).WithMessage("Permission Denied")
}

// IsPermissionDenied 判断是否是权限不足的错误
func IsPermissionDenied(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusForbidden &&
		err.Message == "Permission Denied"
}

// OperationFailed 表示操作失败
func OperationFailed(reason string) *Error {
	return New(http.StatusInternalServerError, reason).WithMessage("Operation Failed")
}

// IsOperationFailed 判断是否是操作失败的错误
func IsOperationFailed(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusInternalServerError &&
		err.Message == "Operation Failed"
}

// Unauthenticated 表示请求未经过身份验证
// 典型场景：API 需要认证但未提供有效凭证
func Unauthenticated(reason string) *Error {
	return New(http.StatusUnauthorized, reason).WithMessage("Unauthenticated")
}

// IsUnauthenticated 判断是否是未认证错误
func IsUnauthenticated(err *Error) bool {
	return err != nil &&
		err.Code == http.StatusUnauthorized &&
		err.Message == "Unauthenticated"
}
