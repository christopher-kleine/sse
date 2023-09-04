package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/christopher-kleine/sse"
)

func (app *App) leave(r *http.Request, s *sse.Session) {
	name := r.FormValue("name")
	room := strings.ToLower(r.FormValue("room"))

	s.Set("Name", name)
	s.Set("Room", room)

	ev := &sse.Event{}
	values := map[string]string{
		"Name": name,
		"Room": room,
		"Time": time.Now().UTC().Format(time.RFC3339),
	}
	t.ExecuteTemplate(ev, "leave.html", values)
	app.hub.FilteredPublish(ev, func(s *sse.Session) bool {
		return s.Get("Room") == room
	})

	var lobby = map[string][]string{}
	for _, session := range app.hub.Sessions() {
		if r := session.Get("Room"); r == room {
			lobby["Lobby"] = append(lobby["Lobby"], session.Get("Name").(string))
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
