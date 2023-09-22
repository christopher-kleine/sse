package main

import (
	"embed"
	"net/http"
	"time"

	"github.com/christopher-kleine/sse"
)

//go:embed index.html
var files embed.FS

func main() {
	hub := sse.New().
		WithHeader("X-DEMO", "demo").
		WithAllowOrigin("chris.isst-gerne.pizza")
	go Ticker(hub)

	http.Handle("/", http.FileServer(http.FS(files)))
	http.Handle("/api/sse", hub)

	http.ListenAndServe(":8080", nil)
}

func Ticker(hub *sse.Hub) {
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		ev := &sse.Event{
			Data: time.Now().UTC().Format(time.RFC3339),
		}

		hub.Publish(ev)
	}
}
