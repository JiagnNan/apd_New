package base

import (
	"github.com/JiagnNan/apd_New/commonlib/dbx"
	"github.com/JiagnNan/apd_New/commonlib/redisx"
	"github.com/JiagnNan/apd_New/server/dao/user"
	"golang.org/x/net/context"
	"testing"
)

func TestSupplementManual(t *testing.T) {
	client, err := redisx.NewRedis(redisx.Config{Addr: "127.0.0.1:16381", DB: 1})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	dbClient, err := dbx.NewMysql(dbx.Config{
		DataSource: "root:root@tcp(127.0.0.1:3306)/apd_test?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai",
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := NewBaseSvc(client, user.NewUserDao(dbClient)).SupplementUserCodes(context.Background()); err != nil {
		t.Fatal(err)
	}
}
