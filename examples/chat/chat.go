package main

import (
	"log"
	"net/http"
)

func (app *App) chat(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Location", "/chat")
		http.SetCookie(w, &http.Cookie{
			Name:  "chat-name",
			Path:  "/",
			Value: r.PostFormValue("name"),
		})
		http.SetCookie(w, &http.Cookie{
			Name:  "chat-room",
			Path:  "/",
			Value: r.PostFormValue("room"),
		})

		w.WriteHeader(http.StatusSeeOther)
		return
	}

	var (
		name, nameErr = r.Cookie("chat-name")
		room, roomErr = r.Cookie("chat-room")
	)

	if nameErr != nil || roomErr != nil {
		log.Println(nameErr, roomErr)
		return
	}

	values := map[string]string{
		"Name": name.Value,
		"Room": room.Value,
	}

	t.ExecuteTemplate(w, "chat.html", values)
}
