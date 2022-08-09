package net

type HandlerFunc func()

type group struct {
	prefix     string // 根据前缀分组	account login||logout
	handlerMap map[string]HandlerFunc
}

type Router struct {
	group []*group
}
