package user

import (
	model "github.com/JiagnNan/apd_New/server/model/user"
	"golang.org/x/net/context"
)

type UserRepository interface {
	GetUserIdByPhoneOrEmail(ctx context.Context, phoneNumber string) (*model.UserInfo, error)
}
