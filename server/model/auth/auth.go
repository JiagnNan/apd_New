package auth

type AuthCode struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Number      string `gorm:"column:number;"`
	Code        string `gorm:"column:code;"`
	CodeType    uint8  `gorm:"column:code_type;"`
	SendType    int64  `gorm:"column:send_type;"`
	ExpiresTime int64  `gorm:"column:expires_time;"`
	Deleted     uint8  `gorm:"column:deleted;"`
}

func (AuthCode) TableName() string {
	return "auth_code"
}
