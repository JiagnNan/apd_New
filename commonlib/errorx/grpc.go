package errorx

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor 把 RPC 业务错误转换为带结构化详情的 gRPC status。
func UnaryServerInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		return nil, ToGRPC(err)
	}
	return resp, nil
}

// ToGRPC 将任意服务端错误收口为可被 API 精确还原的结构化 gRPC 错误。
func ToGRPC(err error) error {
	if err == nil {
		return nil
	}

	if st, ok := status.FromError(err); ok {
		if _, structured := fromGRPCDetails(st); structured {
			return err
		}
	}

	if rpcErr, ok := FromError(err); ok {
		return grpcError(rpcErr)
	}
	if errors.Is(err, context.Canceled) {
		return grpcError(ERR_REQUEST_CANCELED)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return grpcError(ERR_REQUEST_TIMEOUT)
	}

	if st, ok := status.FromError(err); ok {
		return grpcError(errorFromGRPCCode(st.Code()))
	}
	return grpcError(ERR_SERVER_INTERNAL)
}

// grpcError 把业务错误写入标准 ErrorInfo，避免依赖 gRPC message 反推业务码。
func grpcError(rpcErr *RpcError) error {
	info := &errdetails.ErrorInfo{
		Reason: rpcErr.MessageKey,
		Domain: errorDomain,
		Metadata: map[string]string{
			metadataHTTPStatus:   strconv.Itoa(rpcErr.HTTPStatus()),
			metadataBusinessCode: strconv.FormatInt(int64(rpcErr.ErrorCode), 10),
		},
	}

	st := status.New(grpcCodeFromHTTPStatus(rpcErr.HTTPStatus()), rpcErr.MessageKey)
	withDetails, err := st.WithDetails(info)
	if err != nil {
		return status.Error(codes.Internal, ERR_SERVER_INTERNAL.MessageKey)
	}
	return withDetails.Err()
}

// grpcCodeFromHTTPStatus 保留 gRPC 语义，同时让 HTTP 状态由 ErrorInfo 独立传递。
func grpcCodeFromHTTPStatus(httpStatus int) codes.Code {
	switch httpStatus {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusRequestTimeout:
		return codes.Canceled
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	default:
		return codes.Internal
	}
}

// errorFromGRPCCode 为第三方或旧 RPC 错误补齐项目业务码。
func errorFromGRPCCode(code codes.Code) *RpcError {
	switch code {
	case codes.InvalidArgument:
		return ERR_PARAM_INVALID
	case codes.Unauthenticated:
		return ERR_AUTH_KEY_INVALID
	case codes.PermissionDenied:
		return ERR_PERMISSION_DENIED
	case codes.NotFound:
		return ERR_NOT_FOUND
	case codes.Canceled:
		return ERR_REQUEST_CANCELED
	case codes.AlreadyExists, codes.FailedPrecondition, codes.Aborted:
		return NewRpcError(http.StatusConflict, ERR_BAD_REQUEST.ErrorCode, ERR_BAD_REQUEST.MessageKey)
	case codes.ResourceExhausted:
		return ERR_TOO_MANY_REQUESTS
	case codes.DeadlineExceeded:
		return ERR_REQUEST_TIMEOUT
	default:
		return ERR_SERVER_INTERNAL
	}
}
