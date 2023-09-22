package main

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/christopher-kleine/sse"
)

//go:embed static/*
var static embed.FS
var staticFiles, _ = fs.Sub(static, "static")

//go:embed templates/*
var templates embed.FS
var templateFiles, _ = fs.Sub(templates, "templates")
var t, _ = template.ParseFS(templateFiles, "*.html")

func main() {
	app := &App{}
	app.hub = sse.New().
		OnConnect(app.join).
		OnDisconnect(app.leave)

	http.Handle("/", http.FileServer(http.FS(staticFiles)))
	http.HandleFunc("/chat", app.chat)
	http.Handle("/api/sse", app.hub)
	http.HandleFunc("/api/send", app.send)

	http.ListenAndServe(":8080", nil)
}
