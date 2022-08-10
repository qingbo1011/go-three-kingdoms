package api

import (
	"go-three-kingdoms/constant"
	"go-three-kingdoms/db/mysql"
	"go-three-kingdoms/net"
	"go-three-kingdoms/server/login/model"
	"go-three-kingdoms/server/login/proto"

	"github.com/mitchellh/mapstructure"
	logging "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Account struct {
}

var DefaultAccount = &Account{}

func (a *Account) Router(r *net.Router) {
	g := r.Group("account")
	g.AddRouter("login", a.login)
}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	loginReq := &proto.LoginReq{}
	//loginRes := &proto.LoginRsp{}
	// mapstructure参考：https://github.com/mitchellh/mapstructure
	err := mapstructure.Decode(req.Body.Msg, loginReq)
	if err != nil {
		logging.Info("mapstructure.Decode出现错误", err)
	}
	user := model.User{Username: loginReq.Username}
	err = mysql.MysqlDB.Where(&user).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound { // 有没有查出来数据（说明用户不存在）
			rsp.Body.Code = constant.UserNotExist // 返回rsp相应的状态码

		}
		logging.Info("数据库查询出错", err)
		return
	}

	//rsp.Body.Code = 0
	//loginRes := &proto.LoginRsp{
	//	Username: "admin",
	//	Session:  "as",
	//	UId:      1,
	//}
	//rsp.Body.Msg = loginRes
}
