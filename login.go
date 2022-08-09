package main

import (
	"go-three-kingdoms/conf"
	"go-three-kingdoms/net"
)

// http://localhost:8003/api/login
// localhost:8003 服务器  /api/login 路由
// websocket 区别 ：ws://localhost:8003 服务器  发消息 （封装为路由）

func main() {
	server := net.NewServer(conf.LoginHost + ":" + conf.LoginPort)
	server.Start()
}

func init() {
	conf.Init("./conf/config.ini")
}
