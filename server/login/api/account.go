package api

import (
	"go-three-kingdoms/net"
	"go-three-kingdoms/server/login/proto"
)

type Account struct {
}

var DefaultAccount = &Account{}

func (a *Account) Router(r *net.Router) {
	g := r.Group("account")
	g.AddRouter("login", a.login)
}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgRsp) {

	rsp.Body.Code = 0
	loginRes := &proto.LoginRsp{
		Username: "admin",
		Session:  "as",
		UId:      1,
	}
	rsp.Body.Msg = loginRes
}
