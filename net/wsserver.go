package net

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/gorilla/websocket"
	logging "github.com/sirupsen/logrus"
)

// websocket服务
type wsServer struct {
	wsConn       *websocket.Conn
	router       *Router
	outChan      chan *WsMsgRsp
	Seq          int64
	property     map[string]any
	propertyLock sync.RWMutex // 读写锁
}

func NewWsServer(wsConn *websocket.Conn) *wsServer {
	return &wsServer{
		wsConn:   wsConn,
		outChan:  make(chan *WsMsgRsp, 1000),
		Seq:      0,
		property: make(map[string]any),
	}
}

func (w *wsServer) SetRouter(router *Router) {
	w.router = router
}

func (w *wsServer) SetProperty(key string, value any) {
	// 其实Go语言的sync包中提供了一个开箱即用的并发安全版map：sync.Map
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.property[key] = value
}

func (w *wsServer) GetProperty(key string) (any, error) {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()
	if value, ok := w.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (w *wsServer) RemoveProperty(key string) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	delete(w.property, key)
}

func (w *wsServer) Addr() string {
	return w.wsConn.RemoteAddr().String()
}

func (w *wsServer) Push(name string, data any) {
	rsp := &WsMsgRsp{
		Body: &RspBody{
			Name: name,
			Msg:  data,
			Seq:  0,
		},
	}
	w.outChan <- rsp
}

// Start 通道一旦建立,那么收消息和发消息就要一直监听
func (w *wsServer) Start() {
	// 启动读写数据的处理逻辑
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

func (w *wsServer) readMsgLoop() {
	// 先读到客户端发送过来的数据，进行处理后，再回复消息（经过路由）
	defer func() {
		if err := recover(); err != nil {
			logging.Info(err)
			w.Close()
		}
	}()
	for {
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			logging.Info("接收消息出现错误:", err)
			break
		}
		// 收到消息，解析消息（前端发送过来的消息就是json格式）
		// 1.data解压（unzip）
		// 2.前端传来的消息是加密的消息，需要进行解密
		// 3.data转为body（json反序列化）
		body := &ReqBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			logging.Info("数据格式有误，非法格式:", err)
		} else { // 获取到前端传递的数据了，拿上这些数据去具体的业务进行处理
			req := &WsMsgReq{Conn: w, Body: body}
			rsp := &WsMsgRsp{Body: &RspBody{Name: body.Name, Seq: req.Body.Seq}}

			w.outChan <- rsp
		}
	}
}

func (w *wsServer) writeMsgLoop() {
	for {
		select {
		case msg := <-w.outChan:
			w.Write(msg)
		}
	}
}

func (w *wsServer) Write(msg *WsMsgRsp) {
	//data, err := json.Marshal(msg.Body)
	//if err!=nil {
	//	logging.Info(err)
	//}
	//secretKey, err := w.GetProperty("secretKey")
	//if err == nil { // 有加密
	//	key := secretKey.(string)
	//	// 数据做加密
	//}

}

func (w *wsServer) Close() {
	err := w.wsConn.Close()
	if err != nil {
		logging.Info(err)
	}
}
