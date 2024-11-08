package websocketserver

import (
	"github.com/gorilla/websocket"
	"log"
)

// EchoHandler 实现 MessageHandler 接口，用于处理 echo 消息
type EchoHandler struct{}

func (e *EchoHandler) HandleMessage(conn *websocket.Conn, msg Message) {
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("发送 echo 消息失败: %v", err)
	}
}

// BroadcastHandler 实现 MessageHandler 接口，用于处理 broadcast 消息
type BroadcastHandler struct {
	server *WebSocketServer
}

func (b *BroadcastHandler) HandleMessage(conn *websocket.Conn, msg Message) {
	b.server.broadcast <- msg
}

type UpdateCharacters struct {
	server *WebSocketServer
}

func (u *UpdateCharacters) HandleMessage(conn *websocket.Conn, msg Message) {

}

// 在包加载时自动注册处理器
func init() {
	RegisterHandler(MsgTypeEcho, &EchoHandler{})
	RegisterHandler(MsgTypeBroadcast, &BroadcastHandler{})
}
