package errorx

import (
	"errors"
	"net/http"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

const (
	errorDomain          = "apd.errorx"
	metadataHTTPStatus   = "http_status"
	metadataBusinessCode = "business_code"
)

// RpcError 同时保存协议层状态、业务错误码和稳定消息键。
// RPC 层只传递 MessageKey，最终展示文案由 API 层根据语言环境生成。
type RpcError struct {
	HTTPStatusCode int    `json:"-"`
	ErrorCode      int32  `json:"code"`
	MessageKey     string `json:"msg"`
}

// NewRpcError 创建包含 HTTP 状态、业务码和消息键的项目业务错误。
func NewRpcError(httpStatus int, code int32, messageKey string) *RpcError {
	return &RpcError{
		HTTPStatusCode: normalizeHTTPStatus(httpStatus),
		ErrorCode:      code,
		MessageKey:     messageKey,
	}
}

// Error 返回稳定消息键，满足标准 error 接口且不会在 RPC 层提前本地化。
func (e *RpcError) Error() string {
	if e == nil {
		return ERR_SERVER_INTERNAL.MessageKey
	}
	return e.MessageKey
}

// HTTPStatus 返回协议层状态码，业务 code 不参与 HTTP 状态判断。
func (e *RpcError) HTTPStatus() int {
	if e == nil {
		return http.StatusInternalServerError
	}
	return normalizeHTTPStatus(e.HTTPStatusCode)
}

// FromError 从本地 error 或带 ErrorInfo 详情的 gRPC status 中还原业务错误。
func FromError(err error) (*RpcError, bool) {
	if err == nil {
		return nil, false
	}

	var rpcErr *RpcError
	if errors.As(err, &rpcErr) {
		//rpcErr可能是空指针
		if rpcErr == nil {
			return ERR_SERVER_INTERNAL, true
		}
		return rpcErr, true
	}

	st, ok := status.FromError(err)
	if !ok {
		return nil, false
	}
	if rpcErr, ok = fromGRPCDetails(st); ok {
		return rpcErr, true
	}

	// 兼容尚未升级拦截器的旧 RPC 实例；新链路不再依赖错误文案传递业务码。
	if rpcErr, ok = FromMessage(st.Message()); ok {
		return rpcErr, true
	}
	if st.Code() != 0 {
		return errorFromGRPCCode(st.Code()), true
	}
	return nil, false
}

// FromMessage 根据稳定消息键识别预定义错误，仅用于兼容旧 RPC 错误格式。
func FromMessage(messageKey string) (*RpcError, bool) {
	rpcErr, ok := errorByMessageKey[messageKey]
	return rpcErr, ok
}

// fromGRPCDetails 从标准 Error 中读取 HTTP 状态和业务错误码。
func fromGRPCDetails(st *status.Status) (*RpcError, bool) {
	for _, detail := range st.Details() {
		//类型断言：检查这条详情是否为 ErrorInfo
		info, ok := detail.(*errdetails.ErrorInfo)
		if !ok || info.Domain != errorDomain {
			continue
		}

		httpStatus, err := strconv.Atoi(info.Metadata[metadataHTTPStatus])
		if err != nil {
			return nil, false
		}
		businessCode, err := strconv.ParseInt(info.Metadata[metadataBusinessCode], 10, 32)
		if err != nil || info.Reason == "" {
			return nil, false
		}

		return NewRpcError(httpStatus, int32(businessCode), info.Reason), true
	}
	return nil, false
}

// normalizeHTTPStatus 防止业务错误携带非法协议状态码。
func normalizeHTTPStatus(httpStatus int) int {
	if httpStatus >= 400 && httpStatus <= 599 {
		return httpStatus
	}
	return http.StatusInternalServerError
}
