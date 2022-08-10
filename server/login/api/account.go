package api

import (
	"go-three-kingdoms/net"
	"go-three-kingdoms/server/login/proto"

	"github.com/mitchellh/mapstructure"
	logging "github.com/sirupsen/logrus"
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

	//rsp.Body.Code = 0
	//loginRes := &proto.LoginRsp{
	//	Username: "admin",
	//	Session:  "as",
	//	UId:      1,
	//}
	//rsp.Body.Msg = loginRes
}
