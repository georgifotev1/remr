// Package spa serves a single-page application from an fs.FS.
package spa

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// Handler serves static files and falls back to index.html for client routes.
func Handler(dist fs.FS) http.Handler {
	files := http.FileServer(http.FS(dist))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlPath := path.Clean(r.URL.Path)
		if urlPath == "/api" || strings.HasPrefix(urlPath, "/api/") {
			http.NotFound(w, r)
			return
		}
		if strings.HasPrefix(urlPath, "/assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			files.ServeHTTP(w, r)
			return
		}

		name := strings.TrimPrefix(urlPath, "/")
		if name == "" || name == "." {
			serveIndex(w, r, dist)
			return
		}
		if info, err := fs.Stat(dist, name); err == nil && !info.IsDir() {
			files.ServeHTTP(w, r)
			return
		}
		serveIndex(w, r, dist)
	})
}

func serveIndex(w http.ResponseWriter, r *http.Request, dist fs.FS) {
	index, err := fs.ReadFile(dist, "index.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write(index)
}
