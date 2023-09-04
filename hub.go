package sse

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"sync"
)

type Hub struct {
	sessions map[string]*Session

	mu sync.Mutex

	OnConnect    func(*http.Request, *Session)
	OnDisconnect func(*http.Request, *Session)
}

func New() *Hub {
	return &Hub{
		sessions: make(map[string]*Session, 0),
	}
}

func (h *Hub) Publish(ev *Event) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, session := range h.sessions {
		go session.Send(ev)
	}
}

func (h *Hub) FilteredPublish(ev *Event, fn func(*Session) bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, session := range h.sessions {
		if fn(session) {
			go session.Send(ev)
		}
	}
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "server sent events not supported", http.StatusNotAcceptable)
		return
	}

	session := NewSession()
	id := h.AddSession(session)
	session.ID = id

	go session.ServeHTTP(w, r)

	if h.OnConnect != nil {
		h.OnConnect(r, session)
	}

	<-r.Context().Done()
	h.RemoveSession(id)

	if h.OnDisconnect != nil {
		h.OnDisconnect(r, session)
	}
}

func (h *Hub) AddSession(session *Session) string {
	h.mu.Lock()
	defer h.mu.Unlock()

	buffer := make([]byte, 10)
	_, _ = rand.Read(buffer)
	id := fmt.Sprintf("%x", buffer)

	h.sessions[id] = session

	return id
}

func (h *Hub) RemoveSession(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.sessions, id)
}

func (h *Hub) ConnectionCount() int {
	return len(h.sessions)
}

func (h *Hub) Sessions() SessionSlice {
	h.mu.Lock()
	defer h.mu.Unlock()

	result := make(SessionSlice, len(h.sessions))

	k := 0
	for _, v := range h.sessions {
		result[k] = &Session{
			values: v.values,
			ID:     v.ID,
			Joined: v.Joined,
		}
		k++
	}

	return result
}
