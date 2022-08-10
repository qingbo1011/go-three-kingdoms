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
