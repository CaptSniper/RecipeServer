package main

import (
	"bytes"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	rfp "github.com/CaptSniper/RecipeServer/RFP"
)

const frontendDir = "./react/dist"

//go:embed react/dist
var EmbeddedFiles embed.FS

func StartWebServer() {
	cfg, _ := rfp.LoadConfig()

	apiTarget := "http://localhost:" + strconv.Itoa(cfg.DefaultPort)

	// 🔥 Always use embed path, NOT OS path
	fsys, err := fs.Sub(EmbeddedFiles, "react/dist")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(fsys))

	// --- API REVERSE PROXY ---
	apiURL, err := url.Parse(apiTarget)
	if err != nil {
		log.Fatalf("Invalid API URL: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(apiURL)

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
		r.Host = apiURL.Host
		proxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")

		// Try to open file from embedded FS
		f, err := fsys.Open(path)
		if err == nil {
			defer f.Close()
			if info, _ := f.Stat(); info != nil && !info.IsDir() {
				fileServer.ServeHTTP(w, r)
				return
			}
		}

		// SPA fallback to index.html
		indexBytes, err := fs.ReadFile(fsys, "index.html")
		if err != nil {
			http.Error(w, "index.html not found in embedded FS", http.StatusInternalServerError)
			return
		}

		reader := bytes.NewReader(indexBytes)
		http.ServeContent(w, r, "index.html", time.Now(), reader)
	})

	log.Println("Frontend server on http://localhost:" + strconv.Itoa(cfg.DefaultWebPort))
	log.Println("Proxy route: /api →", apiTarget)

	if err := http.ListenAndServe(":"+strconv.Itoa(cfg.DefaultWebPort), nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}

}
