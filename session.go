package sse

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Session struct {
	values map[string]any
	recv   chan *Event
	ID     string
	Joined time.Time
	mu     sync.Mutex
}

type SessionSlice []*Session

func NewSession() *Session {
	return &Session{
		values: make(map[string]any),
		recv:   make(chan *Event),
		Joined: time.Now(),
	}
}

func (s *Session) Send(ev *Event) {
	s.recv <- ev
}

func (s *Session) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "server sent events not supported", http.StatusNotAcceptable)
		return
	}

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

func (s *Session) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.values[key] = value
}

func (s *Session) Get(key string) any {
	if v, ok := s.values[key]; ok {
		return v
	}

	return nil
}
