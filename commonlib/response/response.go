package response

import "github.com/JiagnNan/apd_New/commonlib/errorx"

// Response 是统一的 HTTP JSON 响应体，Code 表示业务码，不是 HTTP 状态码。
type Response struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// Success 创建成功响应，保留调用方传入的数据。
func Success(data any) Response {
	return Response{
		Code: errorx.OK,
		Msg:  "SUCCESS",
		Data: data,
	}
}

// Failure 创建失败响应，data 字段固定输出为 null。
func Failure(code int32, msg string) Response {
	return Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}
