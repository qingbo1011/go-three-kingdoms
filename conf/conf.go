package conf

import (
	"time"

	"github.com/go-ini/ini"
	logging "github.com/sirupsen/logrus"
)

var (
	LoginHost string
	LoginPort string

	MysqlHost                      string
	MysqlPort                      string
	MysqlUser                      string
	MysqlPassword                  string
	MysqlDataBase                  string
	MysqlCharset                   string
	MysqlLogMode                   int
	MysqlDefaultStringSize         uint
	MysqlDisableDatetimePrecision  bool
	MysqlDontSupportRenameIndex    bool
	MysqlDontSupportRenameColumn   bool
	MysqlSkipInitializeWithVersion bool
	MysqlSingularTable             bool
	MysqlMaxIdleConns              int
	MysqlMaxOpenConns              int
	MysqlConnMaxLifetime           time.Duration
)

func Init(path string) {
	file, err := ini.Load(path)
	if err != nil {
		logging.Info("Fail to parse 'conf/app.ini': ", err)
	}

	loadService(file)
	loadMysql(file)
}

func loadService(file *ini.File) {
	section, err := file.GetSection("server")
	if err != nil {
		logging.Info(err)
	}
	LoginHost = section.Key("LoginHost").String()
	LoginPort = section.Key("LoginPort").String()
}

func loadMysql(file *ini.File) {
	section, err := file.GetSection("mysql")
	if err != nil {
		logging.Info(err)
	}
	MysqlHost = section.Key("MysqlHost").String()
	MysqlPort = section.Key("MysqlPort").String()
	MysqlUser = section.Key("MysqlUser").String()
	MysqlPassword = section.Key("MysqlPassword").String()
	MysqlDataBase = section.Key("MysqlDataBase").String()
	MysqlCharset = section.Key("MysqlCharset").String()
	MysqlLogMode = section.Key("MysqlLogMode").MustInt(3)
	MysqlDefaultStringSize = section.Key("MysqlDefaultStringSize").MustUint(256)
	MysqlDisableDatetimePrecision = section.Key("MysqlDisableDatetimePrecision").MustBool(true)
	MysqlDontSupportRenameIndex = section.Key("MysqlDontSupportRenameIndex").MustBool(true)
	MysqlDontSupportRenameColumn = section.Key("MysqlDontSupportRenameColumn").MustBool(true)
	MysqlSkipInitializeWithVersion = section.Key("MysqlSkipInitializeWithVersion").MustBool(false)
	MysqlSingularTable = section.Key("MysqlSingularTable").MustBool(true)
	MysqlMaxIdleConns = section.Key("MysqlMaxIdleConns").MustInt(20)
	MysqlMaxOpenConns = section.Key("MysqlMaxOpenConns").MustInt(100)
	MysqlConnMaxLifetime = time.Duration(section.Key("MysqlConnMaxLifetime").MustInt(30)) * time.Second
}
