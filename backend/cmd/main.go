package main

import (
	"backend/internal/proxy"
	"backend/internal/websocket"
	"backend/pkg/utils"
	"fmt"
	"net/http"
	"os"
)

func getAttributesHandler(w http.ResponseWriter, r *http.Request) {
	// 设置允许跨域访问
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 如果是预检请求，直接返回
	if r.Method == "OPTIONS" {
		return
	}

	// 定义文件路径
	filePath := "./definition/attributes.json"

	// 读取文件内容
	file, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read attributes.json: %v", err), http.StatusInternalServerError)
		return
	}

	// 设置返回类型为 JSON
	w.Header().Set("Content-Type", "application/json")

	// 返回文件内容
	_, err = w.Write(file)
	if err != nil {
		return
	}
}

func main() {
	utils.Version()
	http.HandleFunc("/attributes", getAttributesHandler)
	websocket.LoadCharacters()
	go websocket.HandleMessages()
	http.HandleFunc("/ws", websocket.HandleConnections)
	proxy.ReverseProxy()
}
