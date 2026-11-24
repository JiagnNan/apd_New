package user

import "time"

type UserInfo struct {
	Id          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Code        int64     `gorm:"column:code;"`
	Phone       *string   `gorm:"column:phone;"`
	Email       *string   `gorm:"column:email;"`
	Birthday    time.Time `gorm:"column:birthday;"`
	Gender      uint8     `gorm:"column:gender;"`
	Nickname    string    `gorm:"column:nickname;"`
	Username    *string   `gorm:"column:username;"`
	CountryCode string    `gorm:"column:country_code;"`
	IsLogoff    uint8     `gorm:"column:is_logoff;not null;default:0"`
	IsAuth      uint8     `gorm:"column:is_auth;"`
	LogoffTime  int64     `gorm:"column:logoff_time;"`
	BlackTime   int64     `gorm:"column:black_time;"` // 被封禁的时间戳；0 表示正常，解封时清零，不自动到期。
	WealthLevel int       `gorm:"column:wealth_level;"`
	CharmLevel  int       `gorm:"column:charm_level;"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

type RegisterInfo struct {
	Id     int64   `gorm:"column:id;primaryKey;autoIncrement"`
	Phone  *string `gorm:"column:phone;"`
	Email  *string `gorm:"column:email;"`
	Status uint8   `gorm:"column:status;not null;default:0"`
}

func (RegisterInfo) TableName() string {
	return "register_info"
}
