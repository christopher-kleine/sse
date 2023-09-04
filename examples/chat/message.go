package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/christopher-kleine/sse"
)

func (app *App) send(w http.ResponseWriter, r *http.Request) {
	var (
		name, nameErr = r.Cookie("chat-name")
		room, roomErr = r.Cookie("chat-room")
	)

	if nameErr != nil || roomErr != nil {
		log.Println(nameErr, roomErr)
		return
	}

	roomStr := strings.ToLower(room.Value)

	msg, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	values := map[string]string{
		"Name": name.Value,
		"Room": roomStr,
		"Time": time.Now().UTC().Format(time.RFC3339),
		"Text": string(msg),
	}

	log.Printf("Message: %s", msg)

	ev := &sse.Event{}
	if err := t.ExecuteTemplate(ev, "message.html", values); err != nil {
		log.Println(err)
	}
	app.hub.FilteredPublish(ev, func(s *sse.Session) bool {
		sr, _ := s.Get("Room").(string)
		log.Printf("SR: %v / %s", sr, room.Value)
		return sr == roomStr
	})
}
