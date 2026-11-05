package errorx

import "net/http"

// 动态错误使用独立业务码段，前端可区分审核限制、重复互动和幂等冲突。
var (
	ERR_MOMENT_PENDING          = NewRpcError(http.StatusConflict, 820001, "MOMENT_PENDING")
	ERR_MOMENT_ALREADY_LIKED    = NewRpcError(http.StatusConflict, 820002, "MOMENT_ALREADY_LIKED")
	ERR_MOMENT_NOT_LIKED        = NewRpcError(http.StatusConflict, 820003, "MOMENT_NOT_LIKED")
	ERR_MOMENT_REQUEST_CONFLICT = NewRpcError(http.StatusConflict, 820004, "MOMENT_REQUEST_CONFLICT")
	ERR_MOMENT_AUDIT_COMPLETED  = NewRpcError(http.StatusConflict, 820005, "MOMENT_AUDIT_COMPLETED")
)
