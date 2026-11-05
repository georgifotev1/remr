package spa

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func TestHandler(t *testing.T) {
	dist := fstest.MapFS{
		"index.html":    {Data: []byte("<html>app</html>")},
		"assets/app.js": {Data: []byte("console.log('app')")},
	}
	h := Handler(dist)

	tests := []struct {
		name      string
		path      string
		wantCode  int
		wantBody  string
		wantCache string
	}{
		{
			name:      "index",
			path:      "/",
			wantCode:  http.StatusOK,
			wantBody:  "<html>app</html>",
			wantCache: "no-cache",
		},
		{
			name:      "client route falls back",
			path:      "/clients/42",
			wantCode:  http.StatusOK,
			wantBody:  "<html>app</html>",
			wantCache: "no-cache",
		},
		{
			name:      "asset is cached",
			path:      "/assets/app.js",
			wantCode:  http.StatusOK,
			wantBody:  "console.log('app')",
			wantCache: "public, max-age=31536000, immutable",
		},
		{
			name:     "unknown api path is not the app",
			path:     "/api/v1/missing",
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", w.Code, tt.wantCode)
			}
			if tt.wantBody != "" && w.Body.String() != tt.wantBody {
				t.Errorf("body = %q, want %q", w.Body.String(), tt.wantBody)
			}
			if tt.wantCache != "" {
				if got := w.Header().Get("Cache-Control"); !strings.Contains(got, tt.wantCache) {
					t.Errorf("Cache-Control = %q, want %q", got, tt.wantCache)
				}
			}
		})
	}
}
