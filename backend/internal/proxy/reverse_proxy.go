package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ReverseProxy handles the reverse proxy functionality
func ReverseProxy() {
	// Parse the React server URL
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
		proxy.ServeHTTP(w, r)
	})

	// Start the HTTP server
	log.Println("Starting reverse proxy on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
