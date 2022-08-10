package mysql

import (
	"go-three-kingdoms/conf"
	loginModel "go-three-kingdoms/server/login/model"
	"strings"

	logging "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var MysqlDB *gorm.DB

func Init() {
	// dsn := "root:123456@tcp(127.0.0.1:3308)/test?charset=utf8mb4&parseTime=True&loc=Local"
	var builder strings.Builder
	s := []string{conf.MysqlUser, ":", conf.MysqlPassword, "@tcp(", conf.MysqlHost, ":", conf.MysqlPort, ")/", conf.MysqlDataBase, "?charset=", conf.MysqlCharset, "&parseTime=True&loc=Local"}
	for _, str := range s {
		builder.WriteString(str)
	}
	dsn := builder.String()
	mysqlLogger := logger.Default.LogMode(logger.LogLevel(conf.MysqlLogMode))
	// mysqlLogger := logger.Default   默认LogLevel为Warn，一般我们开发调试时LogLevel设置为Info

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,                                 // DSN data source name
		DefaultStringSize:         conf.MysqlDefaultStringSize,         // string 类型字段的默认长度
		DisableDatetimePrecision:  conf.MysqlDisableDatetimePrecision,  // 禁用datetime精度，MySQL5.6之前的数据库不支
		DontSupportRenameIndex:    conf.MysqlDontSupportRenameIndex,    // 重命名索引时采用删除并新建的方式，MySQL5.7之前的数据库和MariaDB不支持重命名索引
		DontSupportRenameColumn:   conf.MysqlDontSupportRenameColumn,   // 用 `change` 重命名列，MySQL8之前的数据库和MariaDB不支持重命名列
		SkipInitializeWithVersion: conf.MysqlSkipInitializeWithVersion, // 根据当前MySQL版本自动配置
	}), &gorm.Config{
		Logger: mysqlLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: conf.MysqlSingularTable, // 表名不加s
			//TablePrefix:   "test_",                 // 指定表名前缀为test_
		},
	})
	if err != nil {
		logging.Info(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		logging.Info(err)
	}
	sqlDB.SetMaxIdleConns(conf.MysqlMaxIdleConns)       // 设置空闲连接池中的最大连接数
	sqlDB.SetMaxOpenConns(conf.MysqlMaxOpenConns)       // 设置数据库连接最大打开数
	sqlDB.SetConnMaxLifetime(conf.MysqlConnMaxLifetime) // 设置可重用连接的最长时间

	err = db.AutoMigrate(&loginModel.User{}) // 自动迁移
	if err != nil {
		logging.Info(err)
	}

	MysqlDB = db
}
