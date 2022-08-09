package net

import (
	"net/http"

	"github.com/gorilla/websocket"
	logging "github.com/sirupsen/logrus"
)

type server struct {
	addr   string
	router *Router // 一个服务肯定是要是有多个路由的
}

func NewServer(addr string) *server {
	return &server{
		addr: addr,
	}
}

// Start 启动服务
func (s *server) Start() {
	http.HandleFunc("/", s.wsHandler)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		logging.Fatalln(err)
	}
}

// http升级websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { // 允许所有CORS跨域请求
		return true
	},
}

func (s *server) wsHandler(w http.ResponseWriter, r *http.Request) {
	// 1.http协议升级为websocket协议
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Fatalln("websocket服务连接出错", err)
	}
	logging.Println("websocket服务连接成功！")
	// websocket通道建立之后，不管是客户端还是服务端，都可以收发消息
	// 发消息的时候，把消息当做路由处理。消息是有格式的，要先定义消息的格式
	// 客户端发送消息(比如是这个格式{Name:"account.login"})，服务端收到之后进行解析，判断出要处理登录逻辑
	err = wsConn.WriteMessage(websocket.BinaryMessage, []byte("hello"))
	if err != nil {
		logging.Info(err)
	}

}
