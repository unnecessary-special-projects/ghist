package api

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/coderstone/ghist/internal/store"
)

type Server struct {
	store  *store.Store
	mux    *http.ServeMux
	webFS  fs.FS // embedded frontend files (can be nil in dev mode)
	devMode bool
}

func NewServer(s *store.Store, webFS fs.FS, devMode bool) *Server {
	srv := &Server{
		store:   s,
		mux:     http.NewServeMux(),
		webFS:   webFS,
		devMode: devMode,
	}
	srv.routes()
	return srv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	// API routes
	s.mux.HandleFunc("GET /api/tasks", s.handleListTasks)
	s.mux.HandleFunc("POST /api/tasks", s.handleCreateTask)
	s.mux.HandleFunc("GET /api/tasks/{id}", s.handleGetTask)
	s.mux.HandleFunc("PATCH /api/tasks/{id}", s.handleUpdateTask)
	s.mux.HandleFunc("DELETE /api/tasks/{id}", s.handleDeleteTask)
	s.mux.HandleFunc("GET /api/events", s.handleListEvents)
	s.mux.HandleFunc("GET /api/status", s.handleStatus)

	// Serve frontend (embedded or dev proxy)
	if s.webFS != nil {
		s.mux.HandleFunc("/", s.handleFrontend)
	}
}

func (s *Server) handleFrontend(w http.ResponseWriter, r *http.Request) {
	// Don't serve frontend for API routes
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.NotFound(w, r)
		return
	}

	// Try to serve the file directly
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}
	path = strings.TrimPrefix(path, "/")

	f, err := s.webFS.Open(path)
	if err != nil {
		// SPA fallback: serve index.html for any path
		path = "index.html"
		f, err = s.webFS.Open(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}
	f.Close()

	http.ServeFileFS(w, r, s.webFS, path)
}

// cors wraps a handler to add CORS headers (used in dev mode)
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.devMode {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Handler() http.Handler {
	return s.corsMiddleware(s)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error encoding JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
