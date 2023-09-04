package sse

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Event holds the required informations about a specifc event.
type Event struct {
	// ID - if set - results in `id: <VALUE>\n`
	// This can later be retrieved, should a connection get interrupted.
	ID string

	// Event - if set - results in `event: <VALUE>\n`
	// This can be used to trigger a custom event on the client.
	Event string

	// Data - if set - results in `data: <LINE>\n`
	// It can contain multiple lines, such as JSON or HTML.
	Data string

	// Retry is the retry in Millisconds. A value of 0 disanbles this field.
	// If a connection is lost, the client should wait this time before trying to reconnect.
	Retry time.Duration

	buffer []byte
	mu     sync.Mutex
}

// String turns all fields of an Event to valid representation of a SSE chunk.
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

// Write can be used for HTML Templates among other things.
// It appends all writes to the data.
//
// NOTE: Using Write() will replace values set in the Data-field.
func (ev *Event) Write(b []byte) (int, error) {
	ev.mu.Lock()
	defer ev.mu.Unlock()

	ev.buffer = append(ev.buffer, b...)

	return len(b), nil
}

// JSONData takes an object and turns it into a JSON object.
func (ev *Event) JSONData(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	ev.Data = string(data)
	return nil
}
