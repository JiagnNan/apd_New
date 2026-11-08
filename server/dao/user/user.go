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
		query = query.Where("phone = ?", phoneNumber)
	} else {
		query = query.Where("email = ?", phoneNumber)
	}
}
