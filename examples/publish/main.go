package main

import (
	_ "embed"
	"net/http"
	"time"

	"github.com/christopher-kleine/sse"
)

//go:embed index.html
var indexHTML []byte

func main() {
	hub := sse.New()

	go func() {
		ticker := time.NewTicker(10 * time.Second)

		for range ticker.C {
			ev := &sse.Event{
				Data: `
				<p>
					<b>Bold</b> and <i>italic!</i>
				</p>
				`,
			}

			hub.Publish(ev)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(indexHTML)
	})
	http.Handle("/api/sse", hub)

	http.ListenAndServe(":8080", nil)
}
