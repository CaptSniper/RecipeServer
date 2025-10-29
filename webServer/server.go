package webServer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	frontendDir = "./react/dist"
	apiTarget   = "http://localhost:8000"
)

func StartServer() {
	distPath, err := filepath.Abs(frontendDir)
	if err != nil {
		log.Fatalf("Failed to resolve dist directory: %v", err)
	}

	// API Proxy setup
	apiURL, err := url.Parse(apiTarget)
	if err != nil {
		log.Fatalf("Invalid API URL: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(apiURL)

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		// --- CORS headers ---
		w.Header().Set("Access-Control-Allow-Origin", "*") // or restrict to your React URL
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Strip /api prefix so API server sees correct path
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
		r.Host = apiURL.Host

		proxy.ServeHTTP(w, r)
	})

	// Serve static files & SPA fallback
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		if requestPath == "/" {
			requestPath = "/index.html"
		}

		fullPath := filepath.Join(distPath, requestPath)

		// If file exists, serve it
		if stat, err := os.Stat(fullPath); err == nil && !stat.IsDir() {
			http.ServeFile(w, r, fullPath)
			return
		}

		// Otherwise serve SPA entrypoint
		http.ServeFile(w, r, filepath.Join(distPath, "index.html"))
	})

	log.Println("Frontend server running at http://localhost:3000")
	log.Println("Proxy: /api →", apiTarget)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
