package api

import (
	"go-three-kingdoms/constant"
	"go-three-kingdoms/db/mysql"
	"go-three-kingdoms/net"
	"go-three-kingdoms/server/login/model"
	"go-three-kingdoms/server/login/proto"
	"go-three-kingdoms/util"
	"time"

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
	token, err := util.GenerateToken(user.Uid)
	if err != nil {
		logging.Info("GenerateToken出错", err)
	}
	rsp.Body.Code = constant.OK
	loginRes := &proto.LoginRsp{
		Username: user.Username,
		Session:  token,
		UId:      user.Uid,
	}
	rsp.Body.Msg = loginRes

	// 保存用户登录记录（涉及到login_history表）
	loginHistory := model.LoginHistory{
		Uid:      user.Uid,
		State:    constant.Login,
		Ctime:    time.Now(),
		Ip:       loginReq.Ip,
		Hardware: loginReq.Hardware,
	}
	err = mysql.MysqlDB.Create(&loginHistory).Error
	if err != nil {
		logging.Info("存储登录信息出现错误", err)
	}

	// 最后一次登录的状态记录（涉及到login_last表）
	var lastLogin model.LoginLast
	err = mysql.MysqlDB.Where("uid = ?", user.Uid).First(&lastLogin).Error
	if err == gorm.ErrRecordNotFound { // 表中已有数据，需要更新数据
		lastLogin.IsLogout = 0
		lastLogin.Ip = loginReq.Ip
		lastLogin.LoginTime = time.Now()
		lastLogin.Session = token
		lastLogin.Hardware = loginReq.Hardware
		err := mysql.MysqlDB.Model(&lastLogin).Updates(lastLogin).Error
		if err != nil {
			logging.Info("MysqlDB.Model(&lastLogin).Updates(lastLogin)出错", err)
		}
	} else if err == nil { // 没有数据，要插入新数据
		lastLogin.IsLogout = 0
		lastLogin.Ip = loginReq.Ip
		lastLogin.LoginTime = time.Now()
		lastLogin.Session = token
		lastLogin.Hardware = loginReq.Hardware
		lastLogin.Uid = user.Uid
		err := mysql.MysqlDB.Create(&lastLogin).Error
		if err != nil {
			logging.Info("MysqlDB.Create(&lastLogin)出错", err)
		}
	} else {
		logging.Info("MysqlDB.Where(\"uid = ?\", user.Uid).First(&lastLogin)出错", err)
	}

	// 缓存此用户和当前的ws连接

}
