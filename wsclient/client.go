package wsclient

import (
	"github.com/gorilla/websocket" //这里使用的是 gorilla 的 websocket 库
	"log"
)

type Message struct {
	MediaTime int   `json:"mediatime"`
	Ts        int64 `json:"ts"`
}

type WsClient struct {
	Client *websocket.Conn
}

func New() *WsClient {
	dialer := websocket.Dialer{}
	connect, _, _ := dialer.Dial("ws://cloud.cunoe.com:8800/ws", nil)
	return &WsClient{Client: connect}
}

func (wsc WsClient) Close() {
	if wsc.Client != nil {
		wsc.Client.Close()
	}
}

func (wsc WsClient) Writer(msg Message) {
	err := wsc.Client.WriteJSON(&msg)
	if err != nil {
		log.Printf("err: %v", err)
	}
}

func (wsc WsClient) Read(statu chan Message) {
	for {
		var msg Message
		err := wsc.Client.ReadJSON(&msg)
		if err != nil {
			log.Printf("err: %v", err)
		}
		statu <- msg
	}

}
