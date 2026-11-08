package auth

type AuthCode struct {
	Id          int64  `gorm:"column:id;autoIncrement:primaryKey"`
	Number      string `gorm:"column:number;"`
	Code        string `gorm:"column:code;"`
	CodeType    uint8  `gorm:"column:code_type;"`
	SendType    int64  `gorm:"column:send_type;"`
	ExpiresTime int64  `gorm:"column:expires_time;"`
	Deleted     uint8  `gorm:"column:deleted;"`
}
