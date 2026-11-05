package dbx

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// NewMysql 创建 GORM MySQL 连接，并配置连接池和慢 SQL 日志。
func NewMysql(c Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(c.dsn()), &gorm.Config{
		Logger: gormlogger.New(
			//log.New:Go 标准库的日志构造函数
			//Stdout:输出到终端标准输出,prefix:日志前缀,flag:记录日期和时间
			//类似效果:GORM: 2026/11/05 20:55:00 ...
			log.New(os.Stdout, "GORM: ", log.LstdFlags),
			gormlogger.Config{
				//慢阈值
				SlowThreshold:             c.slowThreshold(),
				LogLevel:                  gormlogger.Warn,
				Colorful:                  false, //关闭 GORM 日志的 ANSI 颜色输出
				IgnoreRecordNotFoundError: true,
			},
		),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if c.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	}
	if c.MaxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	}
	sqlDB.SetConnMaxLifetime(c.maxLifetime())
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// MustNewMysql 创建强依赖 MySQL 连接，失败时直接 panic，适用于服务启动阶段。
func MustNewMysql(c Config) *gorm.DB {
	db, err := NewMysql(c)
	if err != nil {
		panic(err)
	}
	return db
}
