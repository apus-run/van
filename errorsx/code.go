package errorsx

var (
	// 2xx
	OK = Success("OK")

	// 4xx
	ErrBadRequest      = BadRequest("BadRequest")
	ErrUnauthorized    = Unauthorized("Unauthorized")
	ErrForbidden       = Forbidden("Forbidden")
	ErrNotFound        = NotFound("NotFound")
	ErrConflict        = Conflict("Conflict")
	ErrPageNotFound    = PageNotFound("NotFound")
	ErrTooManyRequests = TooManyRequests("TooManyRequests")
	ErrClientClosed    = ClientClosed("ClientClosed")

	// 身份验证相关
	ErrTokenInvalid    = TokenInvalid("TokenInvalid")
	ErrTokenExpired    = TokenExpired("TokenExpired")
	ErrTokenInvalidSig = TokenInvalidSignature("TokenInvalidSignature")
	ErrUnauthenticated = Unauthenticated("Unauthenticated")

	// 参数校验
	ErrInvalidParams = InvalidParams("InvalidParams")
	ErrInvalidArgs   = InvalidArguments("InvalidArguments")
	ErrBindError     = BindError("BindError")

	// 5xx
	ErrInternalServer     = InternalServer("InternalServer")
	ErrServiceUnavailable = ServiceUnavailable("ServiceUnavailable")
	ErrGatewayTimeout     = GatewayTimeout("GatewayTimeout")
	ErrPanicError         = PanicError("Panic")

	// 数据库相关
	ErrDBReadError        = DBReadError("DBRead")
	ErrDBWriteError       = DBWriteError("DBWrite")
	ErrDBTransactionError = DBTransactionError("DBTransaction")

	// 业务逻辑
	ErrPermissionDenied = PermissionDenied("PermissionDenied")
	ErrOperationFailed  = OperationFailed("OperationFailed")
)
