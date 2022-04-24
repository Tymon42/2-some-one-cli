package wsclient

import (
	"fmt"
	"github.com/gorilla/websocket" //这里使用的是 gorilla 的 websocket 库
	"log"
	"time"
)

type Message struct {
	Username string `json:"username"`
	Code     string `json:"code"`
	Data     string `json:"data"`
	Ts       int64  `json:"ts"`
}

func wsClient() {
	//创建一个拨号器，也可以用默认的 websocket.DefaultDialer
	dialer := websocket.Dialer{}
	//向服务器发送连接请求，websocket 统一使用 ws://，默认端口和http一样都是80
	connect, _, err := dialer.Dial("ws://cloud.cunoe.com:8800/ws", nil)
	if nil != err {
		log.Println(err)
		return
	}
	//离开作用域关闭连接，go 的常规操作
	defer connect.Close()

	//定时向客户端发送数据
	go tickWriter(connect)

	var msg Message
	//启动数据读取循环，读取客户端发送来的数据
	for {
		//从 websocket 中读取数据
		//messageType 消息类型，websocket 标准
		//messageData 消息数据
		err := connect.ReadJSON(&msg)
		if nil != err {
			log.Println(err)
			break
		}
		//switch messageType {
		//case websocket.TextMessage: //文本数据
		//	fmt.Println(string(messageData))
		//case websocket.BinaryMessage: //二进制数据
		//	fmt.Println(messageData)
		//case websocket.CloseMessage: //关闭
		//case websocket.PingMessage: //Ping
		//case websocket.PongMessage: //Pong
		//default:
		//
		//}
		switch msg.Code {
		case "1":
			//播放函数
		case "2":
			//暂停函数
		case "3":
			//同步函数 其中message应是同步的时间
			fmt.Println(msg.Times)
		}
		//fmt.Printf("recv: %s %s %s %d", msg.Username, msg.Code, msg.Data, msg.Ts)
	}
}

func tickWriter(connect *websocket.Conn) {
	var username string
	var code string
	var data string
	for {
		//向客户端发送类型为文本的数据
		fmt.Println("Please enter your code: // username code data")

		fmt.Scanf("%s %s %s", &username, &code, &data)
		ts := time.Now().UnixNano() / 1e6
		msg := &Message{username, code, data, ts}
		err := connect.WriteJSON(msg)
		if err != nil {
			log.Printf("error: %v", err)
		}
		time.Sleep(1)
	}
}
