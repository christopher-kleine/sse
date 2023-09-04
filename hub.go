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

	OnConnect    func(*Session)
	OnDisconnect func(*Session)
}

// New creates a new SSE-Hub.
func New() *Hub {
	return &Hub{
		sessions: make(map[string]*Session, 0),
	}
}

// Publish let's you publish an event to all connected sessions.
// If you want to send it only to sessions with certain criteria, consider FilteredPublish.
func (h *Hub) Publish(ev *Event) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, session := range h.sessions {
		go session.Send(ev)
	}
}

// FilteredPublish works almost the same as Publish. But it let's you specify a function
// that filters only wanted sessions.
func (h *Hub) FilteredPublish(ev *Event, fn func(*Session) bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, session := range h.sessions {
		if fn(session) {
			go session.Send(ev)
		}
	}
}

// ServeHTTP accepts new SSE connections and adds them to the Session-Pool.
// OnConnect and OnDisconnect are triggered by this function.
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "server sent events not supported", http.StatusNotAcceptable)
		return
	}

	session := NewSession()
	session.Request = r
	id := h.addSession(session)
	session.ID = id

	go session.ServeHTTP(w, r)

	if h.OnConnect != nil {
		h.OnConnect(session)
	}

	<-r.Context().Done()
	h.removeSession(id)

	if h.OnDisconnect != nil {
		h.OnDisconnect(session)
	}
}

func (h *Hub) addSession(session *Session) string {
	h.mu.Lock()
	defer h.mu.Unlock()

	buffer := make([]byte, 10)
	_, _ = rand.Read(buffer)
	id := fmt.Sprintf("%x", buffer)

	h.sessions[id] = session

	return id
}

func (h *Hub) removeSession(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.sessions, id)
}

// ConnectionCouunt returns the currently active sessions/connections
func (h *Hub) ConnectionCount() int {
	return len(h.sessions)
}

// Sessions returns a copy of the current sessions.
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
