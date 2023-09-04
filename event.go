package sse

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Event struct {
	ID    string
	Event string
	Data  string
	Retry time.Duration

	buffer []byte
	mu     sync.Mutex
}

func (ev *Event) String() string {
	result := ""
	if ev.Event != "" {
		result += fmt.Sprintf("event: %s\n", ev.Event)
	}

	if len(ev.buffer) > 0 {
		ev.Data = string(ev.buffer)
	}

	lines := strings.Split(strings.TrimSpace(ev.Data), "\n")
	for _, line := range lines {
		result += fmt.Sprintf("data: %s\n", line)
	}

	if ev.Retry > 0 {
		result += fmt.Sprintf("retry: %d\n", ev.Retry.Milliseconds())
	}
	if ev.ID != "" {
		result += fmt.Sprintf("id: %s\n", ev.ID)
	}

	result += "\n"

	return result
}

func (ev *Event) Write(b []byte) (int, error) {
	ev.mu.Lock()
	defer ev.mu.Unlock()

	ev.buffer = append(ev.buffer, b...)

	return len(b), nil
}

func (ev *Event) JSONData(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	ev.Data = string(data)
	return nil
}
