package dao

import (
	model "github.com/JiagnNan/apd_New/server/model/auth"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type AuthDao struct {
	db *gorm.DB
}

// NewAuthDao 创建认证 DAO。
func NewAuthDao(db *gorm.DB) *AuthDao {
	return &AuthDao{db: db}
}

// DB 将请求上下文绑定到 GORM。
func (d *AuthDao) DB(ctx context.Context) *gorm.DB {
	return d.db.WithContext(ctx)
}

func (d *AuthDao) GetAuthCode(ctx context.Context, phoneNumber string, codeType uint8) (*model.AuthCode, error) {
	var v model.AuthCode
	err := d.DB(ctx).
		Where("phone_number = ?", phoneNumber).
		Where("code_type = ?", codeType).
		Take(&v).Error
	return &v, err
}
func (d *AuthDao) SaveAuthCode(ctx context.Context, authCode *model.AuthCode) error {
	return d.DB(ctx).Save(authCode).Error
}
func (d *AuthDao) DeleteAuthCode(ctx context.Context, phoneNumber string, codeType uint8) error {
	err := d.DB(ctx).Table("auth_code").
		Where("phone_number = ? and deleted = 0", phoneNumber).
		Where("code_type = ?", codeType).
		Update("deleted", 1).Error
	return err
}
