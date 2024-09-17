package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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
	log.Println("version 0.6")
	// Parse the React server URL (IPv4 127.0.0.1)
	reactURL, err := url.Parse("http://frontend:3000")
	if err != nil {
		log.Fatalf("Failed to parse React server URL: %v", err)
	}

	// Create a reverse proxy with the custom transport
	proxy := httputil.NewSingleHostReverseProxy(reactURL)

	// Custom error handler for the proxy
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	}

	// Handle incoming requests and forward them to the React server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Forward the request to the React server
		proxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/attributes", getAttributesHandler)

	// Start the proxy server on port 8080
	log.Println("Starting proxy server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
