package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户信息表
type User struct {
	Uid      uint      `gorm:"column:uid;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"uid"`
	Username string    `gorm:"column:username;type:varchar(20);comment:用户名;NOT NULL" json:"username"`
	Passwd   string    `gorm:"column:passwd;comment:加密后的密文;NOT NULL" json:"passwd"`
	Status   uint8     `gorm:"column:status;type:tinyint(4) unsigned;default:0;comment:用户账号状态。0-默认；1-冻结；2-停号;NOT NULL" json:"status"`
	Hardware string    `gorm:"column:hardware;type:varchar(64);comment:hardware;NOT NULL" json:"hardware"`
	Ctime    time.Time `gorm:"column:ctime;type:timestamp;default:2013-03-15 14:38:09;NOT NULL" json:"ctime"`
	Mtime    time.Time `gorm:"column:mtime;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"mtime"`
}

// TableName 可以实现Tabler接口来更改默认表名
func (User) TableName() string {
	return "user"
}

// SetPassword 密码加密
func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	u.Passwd = string(bytes)
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Passwd), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

// LoginHistory 用户登录表
type LoginHistory struct {
	Id       uint      `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Uid      uint      `gorm:"column:uid;type:int(11) unsigned;default:0;comment:用户UID;NOT NULL" json:"uid"`
	State    uint8     `gorm:"column:state;type:tinyint(4) unsigned;default:0;comment:登录状态，0登录，1登出;NOT NULL" json:"state"`
	Ctime    time.Time `gorm:"column:ctime;type:timestamp;default:CURRENT_TIMESTAMP;comment:登录时间;NOT NULL" json:"ctime"`
	Ip       string    `gorm:"column:ip;type:varchar(31);comment:ip;NOT NULL" json:"ip"`
	Hardware string    `gorm:"column:hardware;type:varchar(64);comment:hardware;NOT NULL" json:"hardware"`
}

// TableName 可以实现Tabler接口来更改默认表名
func (LoginHistory) TableName() string {
	return "login_history"
}

// LoginLast 最后一次用户登录表
type LoginLast struct {
	Id         uint      `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Uid        uint      `gorm:"column:uid;type:int(11) unsigned;default:0;comment:用户UID;NOT NULL" json:"uid"`
	LoginTime  time.Time `gorm:"column:login_time;type:timestamp;comment:登录时间" json:"login_time"`
	LogoutTime time.Time `gorm:"column:logout_time;type:timestamp;comment:登出时间" json:"logout_time"`
	Ip         string    `gorm:"column:ip;type:varchar(50);comment:ip;NOT NULL" json:"ip"`
	IsLogout   uint8     `gorm:"column:is_logout;type:tinyint(4) unsigned;default:0;comment:是否logout,1:logout，0:login;NOT NULL" json:"is_logout"`
	Session    string    `gorm:"column:session;type:varchar(255);comment:会话" json:"session"`
	Hardware   string    `gorm:"column:hardware;type:varchar(64);comment:hardware;NOT NULL" json:"hardware"`
}

// TableName 可以实现Tabler接口来更改默认表名
func (m *LoginLast) TableName() string {
	return "login_last"
}
