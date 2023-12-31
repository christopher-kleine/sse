package main

import (
	"strings"
	"time"

	"github.com/christopher-kleine/sse"
)

func (app *App) join(s *sse.Session) {
	name := s.Request.FormValue("name")
	room := strings.ToLower(s.Request.FormValue("room"))

	s.Set("Name", name)
	s.Set("Room", room)

	ev := &sse.Event{}
	values := map[string]string{
		"Name": name,
		"Room": room,
		"Time": time.Now().UTC().Format(time.RFC3339),
	}
	t.ExecuteTemplate(ev, "join.html", values)
	app.hub.FilteredPublish(ev, func(s *sse.Session) bool {
		return s.Get("Room") == room
	})

	var lobby []string
	for _, session := range app.hub.Sessions() {
		if r := session.Get("Room"); r == room {
			lobby = append(lobby, session.Get("Name").(string))
		}
	}

	ev = &sse.Event{
		Event: "lobby",
	}
	t.ExecuteTemplate(ev, "lobby.html", lobby)
	app.hub.FilteredPublish(ev, func(s *sse.Session) bool {
		return s.Get("Room") == room
	})
}
