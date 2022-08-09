package net

import "strings"

type HandlerFunc func(req *WsMsgReq, rsp *WsMsgRsp)

type group struct {
	prefix     string // 根据前缀分组	account login||logout
	handlerMap map[string]HandlerFunc
}

func (g *group) exec(name string, req *WsMsgReq, rsp *WsMsgRsp) {
	handlerFunc, ok := g.handlerMap[name]
	if ok {
		handlerFunc(req, rsp)
	}
}

func (g *group) AddRouter(name string, handlerFunc HandlerFunc) {
	g.handlerMap[name] = handlerFunc
}

type Router struct {
	group []*group
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Run(req *WsMsgReq, rsp *WsMsgRsp) {
	// req.Body.Name: 路径
	// 以登录业务为例, req.Body.Name=account.login （account为组标识，login为路由标识）
	s := strings.Split(req.Body.Name, ".")
	prefix := ""
	name := ""
	if len(s) == 2 {
		prefix = s[0]
		name = s[1]
	}
	for _, g := range r.group {
		if g.prefix == prefix {
			g.exec(name, req, rsp)
		}
	}
}

func (r *Router) Group(prefix string) *group {
	g := &group{
		prefix:     prefix,
		handlerMap: make(map[string]HandlerFunc),
	}
	r.group = append(r.group, g)
	return g
}
