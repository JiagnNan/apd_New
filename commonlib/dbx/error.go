package dbx

import (
	"errors"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// IsNotFound 判断数据库错误是否为空记录。
func IsNotFound(err error) bool {
	return err != nil && errors.Is(err, gorm.ErrRecordNotFound)
}

// IsDuplicate 判断数据库错误是否为 MySQL 唯一键冲突。
func IsDuplicate(err error) bool {
	var mysqlErr *mysqlDriver.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
