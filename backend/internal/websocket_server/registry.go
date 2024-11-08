package websocketserver

import (
	"github.com/gorilla/websocket"
	"sync"
)

// MessageHandler 是处理器的通用接口
type MessageHandler interface {
	HandleMessage(conn *websocket.Conn, msg Message)
}

// 全局注册表，存储消息类型到处理器的映射
var (
	handlerRegistry = make(map[string]MessageHandler)
	registryMutex   sync.Mutex
)

// RegisterHandler 用于注册新的消息处理器
func RegisterHandler(messageType string, handler MessageHandler) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	handlerRegistry[messageType] = handler
}

// GetHandler 根据消息类型获取对应的处理器
func GetHandler(messageType string) (MessageHandler, bool) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	handler, exists := handlerRegistry[messageType]
	return handler, exists
}
