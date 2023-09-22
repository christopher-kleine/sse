package sse

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Session manages a SSE connection.
// It also contains custom infos set by the developer.
type Session struct {
	Request *http.Request

	headers map[string]string
	values  map[string]any
	recv    chan *Event
	ID      string
	Joined  time.Time
	mu      sync.Mutex
}

type SessionSlice []*Session

// NewSession creates a new Session and makes sure all required members are operational.
func NewSession(headers map[string]string) *Session {
	return &Session{
		values:  make(map[string]any),
		recv:    make(chan *Event),
		headers: headers,
		Joined:  time.Now(),
	}
}

// Send sends an event to the connection.
func (s *Session) Send(ev *Event) {
	s.recv <- ev
}

func (s *Session) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for header, value := range s.headers {
		w.Header().Set(header, value)
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "server sent events not supported", http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return

		case ev := <-s.recv:
			fmt.Fprintf(w, "%s", ev)

			flusher.Flush()
		}
	}
}

// Set sets a custom property.
func (s *Session) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.values[key] = value
}

// Get returns the value of a custom property.
func (s *Session) Get(key string) any {
	if v, ok := s.values[key]; ok {
		return v
	}

	return nil
}
