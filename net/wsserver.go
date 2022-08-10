package net

import (
	"encoding/json"
	"errors"
	"go-three-kingdoms/util"
	"sync"

	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	logging "github.com/sirupsen/logrus"
)

const HandshakeMsg = "handshake"

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
		data, err = util.UnZip(data)
		if err != nil {
			logging.Info("解压数据出错，非法格式：", err)
			continue
		}
		// 2.前端传来的消息是加密的消息，需要进行解密
		secretKey, err := w.GetProperty("secretKey")
		if err == nil { // 有加密
			key := secretKey.(string)
			// 客户端传过来的数据是加密的 需要解密
			d, err := util.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				logging.Info("数据格式有误，解密失败:", err)
				// 出错后，发起握手
				w.Handshake()
			} else {
				data = d
			}
		}
		// 3.data转为body（json反序列化）
		body := &ReqBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			logging.Info("数据格式有误，非法格式:", err)
		} else { // 获取到前端传递的数据了，拿上这些数据去具体的业务进行处理
			req := &WsMsgReq{Conn: w, Body: body}
			rsp := &WsMsgRsp{Body: &RspBody{Name: body.Name, Seq: req.Body.Seq}}
			w.router.Run(req, rsp)
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
	data, err := json.Marshal(msg.Body)
	if err != nil {
		logging.Info(err)
	}
	secretKey, err := w.GetProperty("secretKey")
	if err == nil { // 有加密
		key := secretKey.(string)
		data, err = util.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING) // 数据做加密
		if err != nil {
			logging.Info(err)
		}
	}
	if data, err := util.Zip(data); err == nil { // 压缩
		err := w.wsConn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			logging.Info(err)
		}
	}
}

func (w *wsServer) Close() {
	err := w.wsConn.Close()
	if err != nil {
		logging.Info(err)
	}
}

// Handshake 当游戏客户端发送请求的时候，会先进行握手协议。
// 后端会发送对应的加密key给客户端，这样客户端在发送数据的时候就可以使用此key进行加密处理
func (w *wsServer) Handshake() {
	secretKey := ""
	key, err := w.GetProperty("secretKey")
	if err != nil {
		secretKey = util.RandSeq(16) // 获取secretKey失败，随机生成一个长度为16的字符串
	} else {
		secretKey = key.(string)
	}
	handshake := &Handshake{Key: secretKey}
	body := &RspBody{Name: HandshakeMsg, Msg: handshake}
	data, err := json.Marshal(body)
	if err != nil {
		logging.Info(err)
		return
	}
	if secretKey != "" {
		w.SetProperty("secretKey", secretKey)
	} else {
		w.RemoveProperty("secretKey")
	}
	zipData, err := util.Zip(data)
	if err != nil {
		logging.Info(err)
		return
	}
	err = w.wsConn.WriteMessage(websocket.BinaryMessage, zipData)
	if err != nil {
		logging.Info(err)
	}
}
