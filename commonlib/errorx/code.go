package errorx

import "net/http"

const (
	// OK 表示业务处理成功。
	OK int32 = 200
)

// 系统错误使用 500xxx 码段。
var (
	ERR_SERVER_INTERNAL = NewRpcError(http.StatusInternalServerError, 500001, "SERVER_INTERNAL")
	ERR_INTERNAL        = NewRpcError(http.StatusInternalServerError, 500002, "INTERNAL")
)

// 游戏和陪玩技能错误使用 800xxx 码段。
var (
	ERR_MOBILE_NOT_BOUND         = NewRpcError(http.StatusBadRequest, 800001, "MOBILE_NOT_BOUND")
	ERR_GAME_UNAVAILABLE         = NewRpcError(http.StatusConflict, 800002, "GAME_UNAVAILABLE")
	ERR_SKILL_PRICE_OUT_OF_RANGE = NewRpcError(http.StatusBadRequest, 800003, "SKILL_PRICE_OUT_OF_RANGE")
	ERR_SKILL_APPLICATION_EXISTS = NewRpcError(http.StatusConflict, 800004, "SKILL_APPLICATION_EXISTS")
	ERR_SKILL_STATUS_INVALID     = NewRpcError(http.StatusConflict, 800005, "SKILL_STATUS_INVALID")
	ERR_SKILL_AUDIT_COMPLETED    = NewRpcError(http.StatusConflict, 800006, "SKILL_AUDIT_COMPLETED")
	ERR_SKILL_LAST_LISTED        = NewRpcError(http.StatusConflict, 800007, "SKILL_LAST_LISTED")
)

// 举报错误使用 810xxx 码段。
var (
	ERR_REPORT_SELF          = NewRpcError(http.StatusBadRequest, 810001, "REPORT_SELF")
	ERR_REPORT_ALREADY_TODAY = NewRpcError(http.StatusConflict, 810002, "REPORT_ALREADY_TODAY")
)

// 认证和账号错误使用 100xxx 码段。
var (
	ERR_ACCOUNT_NOT_FOUND         = NewRpcError(http.StatusNotFound, 100001, "ACCOUNT_NOT_FOUND")
	ERR_PASSWORD_INVALID          = NewRpcError(http.StatusUnauthorized, 100002, "PASSWORD_INVALID")
	ERR_CAPTCHA_INVALID           = NewRpcError(http.StatusBadRequest, 100003, "CAPTCHA_INVALID")
	ERR_EMAIL_NUMBER_INVALID      = NewRpcError(http.StatusBadRequest, 100004, "EMAIL_NUMBER_INVALID")
	ERR_PHONE_NUMBER_INVALID      = NewRpcError(http.StatusBadRequest, 100005, "PHONE_NUMBER_INVALID")
	ERR_USER_ALREADY_REGISTERED   = NewRpcError(http.StatusConflict, 100006, "USER_ALREADY_REGISTERED")
	ERR_AUTH_KEY_PERM_EMPTY       = NewRpcError(http.StatusUnauthorized, 100007, "AUTH_KEY_PERM_EMPTY")
	ERR_AUTH_KEY_INVALID          = NewRpcError(http.StatusUnauthorized, 100008, "AUTH_KEY_INVALID")
	ERR_AUTH_KEY_UNREGISTERED     = NewRpcError(http.StatusUnauthorized, 100009, "AUTH_KEY_UNREGISTERED")
	ERR_SESSION_PASSWORD_NEEDED   = NewRpcError(http.StatusUnauthorized, 100010, "SESSION_PASSWORD_NEEDED")
	ERR_SESSION_REVOKED           = NewRpcError(http.StatusUnauthorized, 100011, "SESSION_REVOKED")
	ERR_PHONE_CODE_INVALID        = NewRpcError(http.StatusBadRequest, 100012, "PHONE_CODE_INVALID")
	ERR_ACCOUNT_NOT_REGISTER      = NewRpcError(http.StatusNotFound, 100013, "ACCOUNT_NOT_REGISTER")
	ERR_PHONE_NUMBER_BANNED       = NewRpcError(http.StatusForbidden, 100014, "PHONE_NUMBER_BANNED")
	ERR_CURRENT_CONTACT_NOT_BOUND = NewRpcError(http.StatusBadRequest, 100015, "CURRENT_CONTACT_NOT_BOUND")
	ERR_CONTACT_ALREADY_BOUND     = NewRpcError(http.StatusConflict, 100016, "CONTACT_ALREADY_BOUND")
	ERR_CONTACT_UNCHANGED         = NewRpcError(http.StatusConflict, 100017, "CONTACT_UNCHANGED")
	ERR_CONTACT_VERIFY_INVALID    = NewRpcError(http.StatusBadRequest, 100018, "CONTACT_VERIFY_INVALID")
	ERR_CONTACT_CHANGED           = NewRpcError(http.StatusConflict, 100019, "CONTACT_CHANGED")
	ERR_CONTACT_UPDATE_BUSY       = NewRpcError(http.StatusConflict, 100020, "CONTACT_UPDATE_BUSY")
	ERR_USER_STATUS_ABNORMAL      = NewRpcError(http.StatusBadRequest, 100021, "USER_STATUS_ABNORMAL")
	ERR_ACCOUNT_LOGOFF            = NewRpcError(http.StatusForbidden, 100022, "ACCOUNT_LOGOFF")
	// ERR_USERNAME_ALREADY_EXISTS 表示用户自行设置的用户名已被其他账号占用，由 username 表唯一索引冲突触发。
	ERR_USERNAME_ALREADY_EXISTS = NewRpcError(http.StatusConflict, 100023, "USERNAME_ALREADY_EXISTS")
	// ERR_SMS_COUNTRY_UNAVAILABLE 表示目标国家或地区的短信通道被禁用，该号码无法接收验证码。
	ERR_SMS_COUNTRY_UNAVAILABLE   = NewRpcError(http.StatusForbidden, 100024, "SMS_COUNTRY_UNAVAILABLE")
	ERR_PASSWORD_FORMAT_INVALID   = NewRpcError(http.StatusBadRequest, 100025, "PASSWORD_FORMAT_INVALID")
	ERR_PASSWORD_CONFIRM_MISMATCH = NewRpcError(http.StatusBadRequest, 100026, "PASSWORD_CONFIRM_MISMATCH")
	ERR_PASSWORD_ALREADY_SET      = NewRpcError(http.StatusConflict, 100027, "PASSWORD_ALREADY_SET")
	ERR_PASSWORD_NOT_SET          = NewRpcError(http.StatusConflict, 100028, "PASSWORD_NOT_SET")
	ERR_PASSWORD_UNCHANGED        = NewRpcError(http.StatusConflict, 100029, "PASSWORD_UNCHANGED")
)

// 联系人错误使用 200xxx 码段。
var (
	ERR_FRIEND_ALREADY_EXISTS  = NewRpcError(http.StatusConflict, 200001, "FRIEND_ALREADY_EXISTS")
	ERR_FRIEND_APPLY_NOT_FOUND = NewRpcError(http.StatusNotFound, 200002, "FRIEND_APPLY_NOT_FOUND")
	ERR_FOLLOW_ALREADY_EXISTS  = NewRpcError(http.StatusConflict, 200003, "FOLLOW_ALREADY_EXISTS")
	ERR_FOLLOW_SELF            = NewRpcError(http.StatusBadRequest, 200004, "FOLLOW_SELF")
	ERR_BLOCK_SELF             = NewRpcError(http.StatusBadRequest, 200005, "BLOCK_SELF")
	ERR_BLOCK_ALREADY_EXISTS   = NewRpcError(http.StatusConflict, 200006, "BLOCK_ALREADY_EXISTS")
	ERR_BLOCK_NOT_FOUND        = NewRpcError(http.StatusNotFound, 200007, "BLOCK_NOT_FOUND")
	ERR_BLOCKED_BY_ME          = NewRpcError(http.StatusConflict, 200008, "BLOCKED_BY_ME")
	ERR_BLOCKED_BY_PEER        = NewRpcError(http.StatusConflict, 200009, "BLOCKED_BY_PEER")
)

// 消息和会话错误使用 300xxx 码段。
var (
	ERR_MESSAGE_SEND_FAILED    = NewRpcError(http.StatusBadRequest, 300001, "MESSAGE_SEND_FAILED")
	ERR_CONVERSATION_NOT_FOUND = NewRpcError(http.StatusNotFound, 300002, "CONVERSATION_NOT_FOUND")
	ERR_IM_PEER_UNAVAILABLE    = NewRpcError(http.StatusConflict, 300003, "IM_PEER_UNAVAILABLE")
)

// 文件错误使用 400xxx 码段。
var (
	ERR_FILE_UPLOAD_FAILED = NewRpcError(http.StatusInternalServerError, 400001, "FILE_UPLOAD_FAILED")
)

// 等级配置错误使用 600xxx 码段。
var (
	// ERR_LEVEL_RANGE_INVALID 表示新增等级没有紧接当前最高经验，或发生并发边界冲突。
	ERR_LEVEL_RANGE_INVALID       = NewRpcError(http.StatusConflict, 600001, "LEVEL_RANGE_INVALID")
	ERR_WALLET_NOT_FOUND          = NewRpcError(http.StatusConflict, 600101, "WALLET_NOT_FOUND")
	ERR_WALLET_INSUFFICIENT       = NewRpcError(http.StatusConflict, 600102, "WALLET_INSUFFICIENT")
	ERR_WALLET_OPERATION_CONFLICT = NewRpcError(http.StatusConflict, 600103, "WALLET_OPERATION_CONFLICT")
	ERR_WALLET_AMOUNT_OVERFLOW    = NewRpcError(http.StatusBadRequest, 600104, "WALLET_AMOUNT_OVERFLOW")
)

// 用户资料审核和限制错误使用 700xxx 码段。
var (
	ERR_CONTENT_SENSITIVE            = NewRpcError(http.StatusBadRequest, 700001, "CONTENT_SENSITIVE")
	ERR_ALBUM_LIMIT_EXCEEDED         = NewRpcError(http.StatusBadRequest, 700002, "ALBUM_LIMIT_EXCEEDED")
	ERR_USER_NICKNAME_LENGTH_INVALID = NewRpcError(http.StatusBadRequest, 700003, "USER_NICKNAME_LENGTH_INVALID")
	ERR_USER_NICKNAME_INVALID        = NewRpcError(http.StatusBadRequest, 700004, "USER_NICKNAME_INVALID")
	ERR_AVATAR_INVALID               = NewRpcError(http.StatusBadRequest, 700005, "AVATAR_INVALID")
	ERR_GENDER_INVALID               = NewRpcError(http.StatusBadRequest, 700006, "GENDER_INVALID")
	ERR_BIRTHDAY_INVALID             = NewRpcError(http.StatusBadRequest, 700007, "BIRTHDAY_INVALID")
	ERR_PROVINCE_INVALID             = NewRpcError(http.StatusBadRequest, 700008, "PROVINCE_INVALID")
	ERR_DESCRIPTION_INVALID          = NewRpcError(http.StatusBadRequest, 700009, "DESCRIPTION_INVALID")
	ERR_ALBUNURLS_INVALID            = NewRpcError(http.StatusBadRequest, 700010, "ALBUNURLS_INVALID")
	ERR_LINKS_INVALID                = NewRpcError(http.StatusBadRequest, 700011, "LINKS_INVALID")
	ERR_DEVICES_INVALID              = NewRpcError(http.StatusBadRequest, 700012, "DEVICES_INVALID")
	ERR_TAGS_ID_INVALID              = NewRpcError(http.StatusBadRequest, 700013, "TAGS_ID_INVALID")
	ERR_USER_GAME_LABLE_INVALID      = NewRpcError(http.StatusBadRequest, 700014, "USER_GAME_LABLE_INVALID")
	ERR_USER_ACHIEVEMENT_NULL        = NewRpcError(http.StatusBadRequest, 700015, "USER_ACHIEVEMENT_NULL")
	ERR_USER_AGE_NOT_ALLOWED         = NewRpcError(http.StatusBadRequest, 700017, "USER_AGE_NOT_ALLOWED")
	// ERR_TAG_FAILED_RETRIEVE 表示标签数据读取失败，属于服务端依赖故障，不应报成参数错误。
	ERR_TAG_FAILED_RETRIEVE = NewRpcError(http.StatusInternalServerError, 700016, "TAG_FAILED_RETRIEVE")
)

// 通用协议错误使用 900xxx 码段，避免把 HTTP 状态码直接当成业务码。
var (
	ERR_PARAM_INVALID     = NewRpcError(http.StatusBadRequest, 900001, "PARAM_INVALID")
	ERR_NOT_FOUND         = NewRpcError(http.StatusNotFound, 900002, "NOT_FOUND")
	ERR_TOO_MANY_REQUESTS = NewRpcError(http.StatusTooManyRequests, 900003, "TOO_MANY_REQUESTS")
	ERR_PERMISSION_DENIED = NewRpcError(http.StatusForbidden, 900004, "PERMISSION_DENIED")
	ERR_BAD_REQUEST       = NewRpcError(http.StatusBadRequest, 900005, "BAD_REQUEST")
	ERR_REQUEST_CANCELED  = NewRpcError(http.StatusRequestTimeout, 900006, "REQUEST_CANCELED")
	ERR_REQUEST_TIMEOUT   = NewRpcError(http.StatusGatewayTimeout, 900007, "REQUEST_TIMEOUT")
)

var knownErrors = []*RpcError{
	ERR_SERVER_INTERNAL,
	ERR_INTERNAL,
	ERR_ACCOUNT_NOT_FOUND,
	ERR_PASSWORD_INVALID,
	ERR_CAPTCHA_INVALID,
	ERR_EMAIL_NUMBER_INVALID,
	ERR_PHONE_NUMBER_INVALID,
	ERR_USER_ALREADY_REGISTERED,
	ERR_AUTH_KEY_PERM_EMPTY,
	ERR_AUTH_KEY_INVALID,
	ERR_AUTH_KEY_UNREGISTERED,
	ERR_SESSION_PASSWORD_NEEDED,
	ERR_SESSION_REVOKED,
	ERR_PHONE_CODE_INVALID,
	ERR_ACCOUNT_NOT_REGISTER,
	ERR_PHONE_NUMBER_BANNED,
	ERR_CURRENT_CONTACT_NOT_BOUND,
	ERR_CONTACT_ALREADY_BOUND,
	ERR_CONTACT_UNCHANGED,
	ERR_CONTACT_VERIFY_INVALID,
	ERR_CONTACT_CHANGED,
	ERR_CONTACT_UPDATE_BUSY,
	ERR_USER_STATUS_ABNORMAL,
	ERR_ACCOUNT_LOGOFF,
	ERR_USERNAME_ALREADY_EXISTS,
	ERR_SMS_COUNTRY_UNAVAILABLE,
	ERR_PASSWORD_FORMAT_INVALID,
	ERR_PASSWORD_CONFIRM_MISMATCH,
	ERR_PASSWORD_ALREADY_SET,
	ERR_PASSWORD_NOT_SET,
	ERR_PASSWORD_UNCHANGED,
	ERR_FRIEND_ALREADY_EXISTS,
	ERR_FRIEND_APPLY_NOT_FOUND,
	ERR_FOLLOW_ALREADY_EXISTS,
	ERR_FOLLOW_SELF,
	ERR_BLOCK_SELF,
	ERR_BLOCK_ALREADY_EXISTS,
	ERR_BLOCK_NOT_FOUND,
	ERR_BLOCKED_BY_ME,
	ERR_BLOCKED_BY_PEER,
	ERR_MESSAGE_SEND_FAILED,
	ERR_CONVERSATION_NOT_FOUND,
	ERR_IM_PEER_UNAVAILABLE,
	ERR_FILE_UPLOAD_FAILED,
	ERR_LEVEL_RANGE_INVALID,
	ERR_WALLET_NOT_FOUND,
	ERR_WALLET_INSUFFICIENT,
	ERR_WALLET_OPERATION_CONFLICT,
	ERR_WALLET_AMOUNT_OVERFLOW,
	ERR_CONTENT_SENSITIVE,
	ERR_ALBUM_LIMIT_EXCEEDED,
	ERR_USER_NICKNAME_LENGTH_INVALID,
	ERR_USER_NICKNAME_INVALID,
	ERR_AVATAR_INVALID,
	ERR_GENDER_INVALID,
	ERR_BIRTHDAY_INVALID,
	ERR_PROVINCE_INVALID,
	ERR_DESCRIPTION_INVALID,
	ERR_ALBUNURLS_INVALID,
	ERR_LINKS_INVALID,
	ERR_DEVICES_INVALID,
	ERR_TAGS_ID_INVALID,
	ERR_USER_GAME_LABLE_INVALID,
	ERR_USER_ACHIEVEMENT_NULL,
	ERR_TAG_FAILED_RETRIEVE,
	ERR_USER_AGE_NOT_ALLOWED,
	ERR_MOBILE_NOT_BOUND,
	ERR_GAME_UNAVAILABLE,
	ERR_SKILL_PRICE_OUT_OF_RANGE,
	ERR_SKILL_APPLICATION_EXISTS,
	ERR_SKILL_STATUS_INVALID,
	ERR_SKILL_AUDIT_COMPLETED,
	ERR_SKILL_LAST_LISTED,
	ERR_REPORT_SELF,
	ERR_MOMENT_PENDING,
	ERR_MOMENT_ALREADY_LIKED,
	ERR_MOMENT_NOT_LIKED,
	ERR_MOMENT_REQUEST_CONFLICT,
	ERR_MOMENT_AUDIT_COMPLETED,
	ERR_REPORT_ALREADY_TODAY,
	ERR_PARAM_INVALID,
	ERR_NOT_FOUND,
	ERR_TOO_MANY_REQUESTS,
	ERR_PERMISSION_DENIED,
	ERR_BAD_REQUEST,
	ERR_REQUEST_CANCELED,
	ERR_REQUEST_TIMEOUT,
}

var errorByMessageKey = buildMessageKeyIndex(knownErrors)

// KnownErrors 返回所有已注册错误的副本，便于测试唯一性和生成错误码文档。
func KnownErrors() []*RpcError {
	result := make([]*RpcError, len(knownErrors))
	copy(result, knownErrors)
	return result
}

// FromHTTPStatus 将 go-zero 协议错误转换为带业务码的统一错误。
func FromHTTPStatus(httpStatus int) *RpcError {
	switch httpStatus {
	case http.StatusBadRequest:
		return ERR_PARAM_INVALID
	case http.StatusUnauthorized:
		return ERR_AUTH_KEY_INVALID
	case http.StatusForbidden:
		return ERR_PERMISSION_DENIED
	case http.StatusNotFound:
		return ERR_NOT_FOUND
	case http.StatusRequestTimeout:
		return ERR_REQUEST_CANCELED
	case http.StatusTooManyRequests:
		return ERR_TOO_MANY_REQUESTS
	case http.StatusGatewayTimeout:
		return ERR_REQUEST_TIMEOUT
	default:
		if httpStatus >= http.StatusInternalServerError {
			return NewRpcError(httpStatus, ERR_SERVER_INTERNAL.ErrorCode, ERR_SERVER_INTERNAL.MessageKey)
		}
		return NewRpcError(httpStatus, ERR_BAD_REQUEST.ErrorCode, ERR_BAD_REQUEST.MessageKey)
	}
}

// Msg 根据业务码返回稳定消息键；未知业务码统一回退到系统错误。
func Msg(code int32) string {
	if code == OK {
		return "SUCCESS"
	}
	for _, rpcErr := range knownErrors {
		if rpcErr.ErrorCode == code {
			return rpcErr.MessageKey
		}
	}
	return ERR_SERVER_INTERNAL.MessageKey
}

// buildMessageKeyIndex 为旧 RPC 文案兼容路径建立只读索引。
func buildMessageKeyIndex(errors []*RpcError) map[string]*RpcError {
	result := make(map[string]*RpcError, len(errors))
	for _, rpcErr := range errors {
		result[rpcErr.MessageKey] = rpcErr
	}
	return result
}
