package auth

type CodeType uint8

const (
	CodeTypeUnknown CodeType = iota
	CodeTypeLogin
	CodeTypeRegister
)

type LoginType uint8

const (
	LoginTypeUnknown LoginType = iota
	LoginTypeLogin
	LoginTypeRegister
	LoginOutType
)
