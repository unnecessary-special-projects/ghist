package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// hub manages SSE subscriber channels.
type hub struct {
	mu      sync.Mutex
	clients map[chan struct{}]struct{}
}

func newHub() *hub {
	return &hub{clients: make(map[chan struct{}]struct{})}
}

func (h *hub) subscribe() chan struct{} {
	ch := make(chan struct{}, 1)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *hub) unsubscribe(ch chan struct{}) {
	h.mu.Lock()
	delete(h.clients, ch)
	h.mu.Unlock()
}

func (h *hub) broadcast() {
	h.mu.Lock()
	for ch := range h.clients {
		select {
		case ch <- struct{}{}:
		default: // already pending, drop
		}
	}
	h.mu.Unlock()
}

// watchDirs polls dirs every 500ms and broadcasts when anything changes.
func watchDirs(h *hub, dirs []string) {
	last := dirFingerprint(dirs)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		fp := dirFingerprint(dirs)
		if fp != last {
			last = fp
			h.broadcast()
		}
	}
}

// dirFingerprint returns a string that changes whenever any file in dirs
// is created, deleted, or modified.
func dirFingerprint(dirs []string) string {
	var sb strings.Builder
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			info, err := e.Info()
			if err != nil {
				continue
			}
			fmt.Fprintf(&sb, "%s:%d:%d|", e.Name(), info.Size(), info.ModTime().UnixNano())
		}
	}
	return sb.String()
}

// handleSSE streams server-sent events to the client. It blocks until the
// client disconnects, sending "data: update\n\n" whenever data changes.
func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // disable nginx/proxy buffering

	ch := s.hub.subscribe()
	defer s.hub.unsubscribe(ch)

	for {
		select {
		case <-ch:
			fmt.Fprintf(w, "data: update\n\n")
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
