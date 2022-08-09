package login

import (
	"go-three-kingdoms/net"
	"go-three-kingdoms/server/login/api"
)

var Router = net.NewRouter() //var Router = &net.Router{}

func Init() {
	initRouter()
}

func initRouter() {
	api.DefaultAccount.Router(Router)
}
