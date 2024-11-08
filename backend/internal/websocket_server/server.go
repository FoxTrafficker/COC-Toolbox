package websocketserver

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketServer 管理 WebSocket 连接和消息处理
type WebSocketServer struct {
	clients   map[*websocket.Conn]bool
	upgrader  websocket.Upgrader
	mu        sync.Mutex
	broadcast chan Message
}

// NewWebSocketServer 创建新的 WebSocketServer
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients:   make(map[*websocket.Conn]bool),
		upgrader:  websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		broadcast: make(chan Message),
	}
}

// Start 启动 WebSocket 服务器
func (ws *WebSocketServer) Start(addr string) error {
	http.HandleFunc("/ws", ws.handleConnections)
	go ws.handleMessages()
	log.Printf("WebSocket 服务器启动，监听端口 %s", addr)
	return http.ListenAndServe(addr, nil)
}

// handleConnections 处理新的客户端连接
func (ws *WebSocketServer) handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("无法升级为 WebSocket 连接: %v", err)
		return
	}

	ws.mu.Lock()
	ws.clients[conn] = true
	ws.mu.Unlock()

	defer func() {
		ws.mu.Lock()
		delete(ws.clients, conn)
		ws.mu.Unlock()
		err := conn.Close()
		if err != nil {
			log.Printf("WebSocket 连接在关闭时出现问题: %v", err)
			return
		}
	}()

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("读取消息失败: %v", err)
			break
		}
		ws.routeMessage(conn, msg)
	}
}

// routeMessage 使用注册表动态路由到处理器
func (ws *WebSocketServer) routeMessage(conn *websocket.Conn, msg Message) {
	if handler, exists := GetHandler(msg.Type); exists {
		handler.HandleMessage(conn, msg)
	} else {
		log.Printf("未知的消息类型: %s", msg.Type)
	}
}

// handleMessages 处理广播消息，将其发送给所有客户端
func (ws *WebSocketServer) handleMessages() {
	for {
		msg := <-ws.broadcast
		ws.mu.Lock()
		for client := range ws.clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Printf("广播消息失败: %v", err)
				err := client.Close()
				if err != nil {
					return
				}
				delete(ws.clients, client)
			}
		}
		ws.mu.Unlock()
	}
}
