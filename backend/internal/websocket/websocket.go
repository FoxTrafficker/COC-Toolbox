package websocket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var clients = make(map[*websocket.Conn]bool) // 存储连接的客户端
var broadcast = make(chan Message)           // 广播消息的通道

var mu sync.Mutex          // 用于并发控制
var characters []Character // 存储当前的角色状态

type Character struct {
	Name    string `json:"name"`
	Agility int    `json:"agility"`
	Avatar  string `json:"avatar"`
}

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// 处理 WebSocket 连接

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("升级 WebSocket 连接失败:", err)
		return
	} else {
		log.Printf("与 %s WebSocket 建立", r.RemoteAddr)
	}

	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Printf("与 %s 的WebSocket链接未正常关闭：%s", r.RemoteAddr, err.Error())
		} else {
			log.Printf("与 %s 的WebSocket链接关闭", r.RemoteAddr)
		}
	}(ws)

	// 将新连接加入客户端列表
	clients[ws] = true
	log.Printf("当前加入的链接：")
	for client := range clients {
		log.Printf("%s", client.RemoteAddr().String())
	}

	// 发送当前的角色状态给新连接的客户端
	mu.Lock()
	initialState, _ := json.Marshal(characters)
	mu.Unlock()

	ws.WriteJSON(Message{
		Type:    "INITIAL_STATE",
		Payload: string(initialState),
	})

	for {
		var msg Message
		// 读取消息
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println("读取消息失败:", err)
			delete(clients, ws)
			break
		}

		// 根据消息类型处理
		switch msg.Type {
		case "UPDATE_CHARACTERS":
			var updatedCharacters []Character
			err := json.Unmarshal([]byte(msg.Payload), &updatedCharacters)
			if err == nil {
				mu.Lock()
				characters = updatedCharacters
				mu.Unlock()

				// 将角色状态保存到文件
				saveCharacters()

				// 广播更新给所有客户端
				broadcast <- msg
			}

		case "RESET_CHARACTERS":
			// Sort characters by agility in descending order
			mu.Lock()
			sort.Slice(characters, func(i, j int) bool {
				return characters[i].Agility > characters[j].Agility
			})
			sortedCharacters, _ := json.Marshal(characters)
			mu.Unlock()

			// Save the sorted characters to the file
			saveCharacters()

			// Broadcast the updated sorted characters to all clients
			broadcast <- Message{
				Type:    "UPDATE_CHARACTERS",
				Payload: string(sortedCharacters),
			}
		}
	}
}

func HandleMessages() {
	for {
		select {
		case msg := <-broadcast: // 从广播通道中读取消息
			mu.Lock()
			for client := range clients {
				go func(c *websocket.Conn) {
					// 尝试带有重试机制发送消息
					if err := sendWithRetries(c, msg, 3); err != nil {
						fmt.Println("发送消息失败，移除客户端:", err)
						c.Close()
						mu.Lock()
						delete(clients, c) // 出错时移除客户端
						mu.Unlock()
					}
				}(client)
			}
			mu.Unlock()
		}
	}
}

func sendWithRetries(client *websocket.Conn, msg interface{}, retries int) error {
	for i := 0; i < retries; i++ {
		err := client.WriteJSON(msg)
		if err == nil {
			return nil // 发送成功
		}

		// 判断是否是暂时性错误（网络抖动等）
		if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
			fmt.Println("暂时性错误，重试发送:", err)
			continue // 重试
		}

		// 其他错误或者重试超时
		return err
	}
	return fmt.Errorf("发送消息失败，超过最大重试次数")
}

func saveCharacters() {
	mu.Lock()
	defer mu.Unlock()
	data, _ := json.Marshal(characters)
	_ = ioutil.WriteFile("db/characters.json", data, 0644)
}

func LoadCharacters() {
	file, err := os.Open("db/characters.json")
	if err != nil {
		fmt.Println("无法打开保存的角色文件，初始化为空状态")
		characters = []Character{}
		return
	}
	defer file.Close()

	data, _ := ioutil.ReadAll(file)
	_ = json.Unmarshal(data, &characters)
}
