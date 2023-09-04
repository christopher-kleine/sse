# SSE - Melody inspired Server Sent Events Library

[![go.dev Reference](https://pkg.go.dev/static/frontend/badge/badge.svg)](https://go.dev/github.com/christpher-kleine/sse)

A small library providing easy SSE functionality. The API is inspired by the gread WebSocket library [melody](https://github.com/olahol/melody).

**Note:** This library is mainly for my own use and didn't account for other use-cases. If you need something more stable or a more mature library, check [r3labs/sse](https://github.com/r3labs/sse). Same goes for if you want a SSE Client library. But please feel free to tear the source and I'm open for PRs.

## Install

```
go get github.com/christopher-kleine/sse
```

## How to use it

```go
package main

import (
    "net/http"

    "github.com/christopher-kleine/sse"
)

func main() {
    // Create a new Hub
    hub := sse.New()

    // (Optional): Specify OnConnect and OnDisconnect hooks
    hub.OnConnect = func(session *sse.Session) {}
    hub.OnDisconnect = func(session *sse.Session) {}

    // Specify the endpoint
    http.Handle("/sse", hub)

    // (Optional): Customize the request
    /*
    http.HandleFunc("/sse", func(w http.Response, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
    })
    */

    // Publish some data!
    go func() {
        for {
            hub.Publish(&sse.Event{ Data: "Hello World" })
            time.Sleep(1 * time.Second)
        }
    }()

    http.ListenAndServe(":8080", nil)
}
```

## Filtered/Selected Publish

You can also publish to selected sessions.

```go
// Only sent to users we gave the "villain" role.
ev := &sse.Event{Data: "Hello, Villain. What are your next plans?"}
hub.FilteredPublish(ev, func(session *sse.Session) bool {
    return session.Get("role") != "villain"
})
```

## HTMX / HTML Templates

You can use this library to send HTML templates over SSE, since the `Event` type implements the `io.Writer` Interface:

```go
ev := &sse.Event{}
templates.ExecuteTemplate(ev, "index.html", nil)
hub.Publish(ev)
```