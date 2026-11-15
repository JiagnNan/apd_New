package auth

type CodeType uint8

const (
	CodeTypeUnknown CodeType = iota
	CodeTypeLogin
	CodeTypeRegister
)

type ActionType uint8

const (
	ActionTypeUnknown ActionType = iota
	ActionTypeLogin
	ActionTypeRegister
	ActionTypeType //退出
)

type LoginType uint8

const (
	LoginTypeUnknown LoginType = iota
	PasswordLoginType
	PhoneLoginType
	EmailLoginType
)
