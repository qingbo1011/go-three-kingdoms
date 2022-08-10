package model

import "time"

// LoginHistory 用户登录表
type LoginHistory struct {
	Id       uint      `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Uid      uint      `gorm:"column:uid;type:int(11) unsigned;default:0;comment:用户UID;NOT NULL" json:"uid"`
	State    uint8     `gorm:"column:state;type:tinyint(4) unsigned;default:0;comment:登录状态，0登录，1登出;NOT NULL" json:"state"`
	Ctime    time.Time `gorm:"column:ctime;type:timestamp;default:CURRENT_TIMESTAMP;comment:登录时间;NOT NULL" json:"ctime"`
	Ip       string    `gorm:"column:ip;type:varchar(31);comment:ip;NOT NULL" json:"ip"`
	Hardware string    `gorm:"column:hardware;type:varchar(64);comment:hardware;NOT NULL" json:"hardware"`
}

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

func (m *LoginLast) TableName() string {
	return "login_last"
}
