package base

import (
	"github.com/JiagnNan/apd_New/commonlib/pkg/any32"
	"math/rand"
)

func (b *BaseSvc) MakeRandomUsername(userId int64) string {
	//00 10000003 000
	//前2位随机 中间8位userId 后3位随机数
	RandNumber := rand.Int63n(99) * userId * 1000 * rand.Int63n(999)
	randUsername := any32.DecimalToAny32(RandNumber)
	randUsername = "qy_" + randUsername
	return randUsername
}
