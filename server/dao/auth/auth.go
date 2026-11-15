package auth

import (
	model "github.com/JiagnNan/apd_New/server/model/auth"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"time"
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

func (d *AuthDao) GetAuthCode(ctx context.Context, phoneNumber string, codeType model.CodeType) (*model.AuthCode, error) {
	var v model.AuthCode
	err := d.DB(ctx).
		Where("number = ?", phoneNumber).
		Where("code_type = ?", codeType).
		Take(&v).Error
	return &v, err
}
func (d *AuthDao) SaveAuthCode(ctx context.Context, authCode *model.AuthCode) error {
	return d.DB(ctx).Save(authCode).Error
}
func (d *AuthDao) DeleteAuthCode(ctx context.Context, phoneNumber string, codeType model.CodeType) error {
	err := d.DB(ctx).Table("auth_code").
		Where("number = ? and deleted = 0", phoneNumber).
		Where("code_type = ?", codeType).
		Update("deleted", 1).Error
	return err
}
func (d *AuthDao) ConsumeAuthCode(ctx context.Context, code *model.AuthCode) bool {
	result := d.DB(ctx).Model(&model.AuthCode{}).
		Where("id = ? AND code = ? AND expires_time = ?",
			code.ID, code.Code, code.ExpiresTime).
		Where("deleted = ? AND expires_time > ?", 0, time.Now().Unix()).
		Update("deleted", 1)

	if result.Error != nil {
		return false
	}
	return result.RowsAffected == 1
}

func (d *AuthDao) Login(ctx context.Context, authLogin *model.AuthLogin) error {
	return d.DB(ctx).Create(&authLogin).Error
}
