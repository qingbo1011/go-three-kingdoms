package api

import (
	"fmt"
	"go-three-kingdoms/conf"
	"go-three-kingdoms/db/mysql"
	"go-three-kingdoms/log"
	"go-three-kingdoms/server/login/model"
	"testing"
	"time"

	logging "github.com/sirupsen/logrus"
)

func init() {
	log.Init()
	conf.Init("../../../conf/config.ini")
	mysql.Init()
}

func TestRegister(t *testing.T) {
	user := model.User{
		Uid:      1001,
		Username: "test",
		Status:   0,
		Hardware: "hadrware",
		Ctime:    time.Now(),
		Mtime:    time.Now().Add(time.Hour * 1),
	}
	err := user.SetPassword("1234")
	if err != nil {
		logging.Info("user.SetPassword出现错误", err)
	}
	err = mysql.MysqlDB.Create(&user).Error
	if err != nil {
		logging.Info(err)
	}
}

func TestSelectUserByUserName(t *testing.T) {
	user := model.User{Username: "test"}
	err := mysql.MysqlDB.Where(&user).First(&user).Error
	if err != nil {
		logging.Info(err)
		return
	}
	fmt.Println(user)
}
