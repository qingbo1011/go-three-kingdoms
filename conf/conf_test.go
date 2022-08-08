package conf

import (
	"testing"
)

func init() {
	Init("./config.ini")
}

func TestConf(t *testing.T) {
	//fmt.Println(LoginHost)
	//fmt.Println(MysqlDataBase)
	//fmt.Println(MysqlLogMode)
	//fmt.Println(MysqlHost)
}
