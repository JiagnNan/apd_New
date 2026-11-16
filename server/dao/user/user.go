package user

import (
	model "github.com/JiagnNan/apd_New/server/model/user"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"strings"
)

type UserDao struct {
	db *gorm.DB
}

// NewAuthDao 创建认证 DAO。
func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

// DB 将请求上下文绑定到 GORM。
func (d *UserDao) DB(ctx context.Context) *gorm.DB {
	return d.db.WithContext(ctx)
}

func (d *UserDao) GetUserIdByPhoneOrEmail(ctx context.Context, phoneNumber string) (*model.UserInfo, error) {
	var v model.UserInfo
	query := d.DB(ctx)
	if strings.Contains(phoneNumber, "@") {
		query = query.Where("email = ?", phoneNumber)
	} else {
		query = query.Where("phone = ?", phoneNumber)
	}
	err := query.First(&v).Error

	return &v, err
}
func (d *UserDao) CreateUserInfo(ctx context.Context, v *model.UserInfo) error {
	return d.DB(ctx).Create(&v).Error
}

func (d *UserDao) UpdateUserNickName(ctx context.Context, userId int64, nickname string) error {
	return d.DB(ctx).
		Table("user_info").
		Where("id = ?", userId).
		Update("nickname", nickname).Error
}

func (d *UserDao) InsertRegisterInfo(ctx context.Context, v *model.RegisterInfo) error {
	return d.DB(ctx).Create(&v).Error
}
