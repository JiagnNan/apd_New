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

type AuthLogin struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement"`
	UserId        int64  `gorm:"column:user_id;"`
	Ip            string `gorm:"column:ip;"`
	DeviceId      string `gorm:"column:device_id;"`
	DeviceModel   string `gorm:"column:device_model;"`
	SystemVersion string `gorm:"column:system_version;"`
	AppVersion    string `gorm:"column:app_version;"`
	Platform      string `gorm:"column:platform;"`
	Language      string `gorm:"column:language;"`
	LoginType     uint8  `gorm:"column:login_type;"`
	ActionType    uint8  `gorm:"column:action_type;"`
	ActionTime    int64  `gorm:"column:action_time;"`
	CountryCode   string `gorm:"column:country_code;"`
	Deleted       uint8  `gorm:"column:deleted;"`
}

func (AuthLogin) TableName() string {
	return "auth_login"
}
