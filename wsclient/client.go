package wsclient

import (
	"github.com/gorilla/websocket" //这里使用的是 gorilla 的 websocket 库
	"log"
)

var statu = make(chan Message)

type Message struct {
	Mediatime int   `json:"mediatime"`
	Ts        int64 `json:"ts"`
}

type WsClient struct {
	client *ws.Conn
}

func (wsc WsClient) wsClient() *ws.Conn {
	//创建一个拨号器，也可以用默认的 websocket.DefaultDialer
	dialer := websocket.Dialer{}
	//向服务器发送连接请求，websocket 统一使用 ws://，默认端口和http一样都是80
	var err error
	connect, _, err := dialer.Dial("ws://cloud.cunoe.com:8800/ws", nil)
	if nil != err {
		log.Println(err)
		return nil
	}
	wsc.client := connect
	return wsc.client
	////离开作用域关闭连接，go 的常规操作
	//defer connect.Close()
	//
	//var msg Message
	////启动数据读取循环，读取客户端发送来的数据
	//for {
	//	//从 websocket 中读取数据
	//	//messageType 消息类型，websocket 标准
	//	//messageData 消息数据
	//	err := connect.ReadJSON(&msg)
	//	if nil != err {
	//		log.Println(err)
	//		break
	//	}
	//	//switch messageType {
	//	//case websocket.TextMessage: //文本数据
	//	//	fmt.Println(string(messageData))
	//	//case websocket.BinaryMessage: //二进制数据
	//	//	fmt.Println(messageData)
	//	//case websocket.CloseMessage: //关闭
	//	//case websocket.PingMessage: //Ping
	//	//case websocket.PongMessage: //Pong
	//	//default:
	//	//
	//	//}
	//	staut <- msg
	//	//fmt.Printf("recv: %s %s %s %d", msg.Username, msg.Code, msg.Data, msg.Ts)
	//}
	//return WsClient{client: connect}
}

func (wsc WsClient) close(){
	if wsc.client != nil{
		wsc.client.Close()
	}
}

func (wsc WsClient) Writer(msg Message) {
	err := wsc.WriteJSON(&msg)
	if err != nil {
		log.Printf("err: %v", err)
	}
}

func (wsc WsClient) Read(statu chan Message) {
	for {
		var msg Message
		err := wsc.ReadJSON(&msg)
		if err != nil {
			log.Printf("err: %v", err)
		}
		statu <- msg
	}

}