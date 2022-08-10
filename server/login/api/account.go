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
	loginRes := &proto.LoginRsp{}
	// mapstructure参考：https://github.com/mitchellh/mapstructure
	err := mapstructure.Decode(req.Body.Msg, loginReq)
	if err != nil {
		logging.Info("mapstructure.Decode出现错误", err)
	}
	// 查询用户，判断能否登录成功（涉及到user表）
	user := model.User{Username: loginReq.Username}
	err = mysql.MysqlDB.Where(&user).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound { // 有没有查出来数据（说明用户不存在）
			rsp.Body.Code = constant.UserNotExist // 返回rsp相应的状态码
			return
		}
		logging.Info("数据库查询出错", err)
		return
	}
	// 查到user后，比较密码
	ok, _ := user.CheckPassword(loginReq.Password)
	if !ok { // 密码不正确
		rsp.Body.Code = constant.PwdIncorrect
	}
	// 用户名和密码都没问题，开始签发token
	// jwt A.B.C三部分：A定义加密算法，B定义放入的数据，C部分根据秘钥+A和B生成加密字符串
	//（更详细的介绍参考：https://juejin.cn/post/6844904090690912264）

	//rsp.Body.Code = 0
	//loginRes := &proto.LoginRsp{
	//	Username: "admin",
	//	Session:  "as",
	//	UId:      1,
	//}
	//rsp.Body.Msg = loginRes
}
