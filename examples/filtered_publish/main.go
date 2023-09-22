package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/christopher-kleine/sse"
)

func main() {
	var count int = 0

	hub := sse.New().
		OnConnect(func(session *sse.Session) {
			session.Set("count", count)
			count++
			session.Send(&sse.Event{
				Data: "HELLO WORLD!",
			})
			log.Printf("Connected: IP %v connected with %+v", session.Request.RemoteAddr, session)
		}).
		OnDisconnect(func(session *sse.Session) {
			log.Printf("DisConnected: IP %v connected with %+v", session.Request.RemoteAddr, session)
		})

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		id := 0

		for range ticker.C {
			hub.FilteredPublish(&sse.Event{
				Event: "foo",
				Data:  fmt.Sprintf("bar (%d / %d)", count, hub.ConnectionCount()),
				ID:    fmt.Sprintf("msg-%d", id),
				Retry: 30 * time.Second,
			}, func(session *sse.Session) bool {
				c := session.Get("count").(int)

				return c%2 == 0
			})

			hub.FilteredPublish(&sse.Event{
				Event: "bar",
				Data:  fmt.Sprintf("FU! (%d / %d)", count, hub.ConnectionCount()),
				ID:    fmt.Sprintf("msg-%d", id),
				Retry: 30 * time.Second,
			}, func(session *sse.Session) bool {
				c := session.Get("count").(int)

				return c%2 == 1
			})

			id++
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/api/sse", hub)

	http.ListenAndServe(":8080", nil)
}
