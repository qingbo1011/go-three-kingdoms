package net

type ReqBody struct {
	Seq   int64  `json:"seq"`
	Name  string `json:"name"`
	Msg   any    `json:"msg"`
	Proxy string `json:"proxy"`
}

type RspBody struct {
	Seq  int64  `json:"seq"`
	Name string `json:"name"`
	Code int    `json:"code"`
	Msg  any    `json:"msg"`
}

type WsMsgReq struct {
	Body *ReqBody
	Conn WSConn
}

type WsMsgRsp struct {
	Body *RspBody
}

type Handshake struct {
	Key string `json:"key"`
}

// WSConn request请求会有参数，请求中放参数、取参数的逻辑
type WSConn interface {
	SetProperty(key string, value any)
	GetProperty(key string) (any, error)
	RemoveProperty(key string)
	Addr() string
	Push(name string, data any)
}
