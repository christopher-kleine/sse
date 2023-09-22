# SSE - Melody inspired Server Sent Events Library

[![go.dev Reference](https://pkg.go.dev/static/frontend/badge/badge.svg)](https://pkg.go.dev/github.com/christopher-kleine/sse) ![Test Status](https://github.com/christopher-kleine/sse/actions/workflows/test.yml/badge.svg) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![latest tag](https://img.shields.io/github/v/tag/christopher-kleine/sse)](https://github.com/christopher-kleine/sse/tags)

## Synopsis

A small library providing easy SSE functionality. The API is inspired by the great WebSocket library [melody](https://github.com/olahol/melody).

**Note:** This library is mainly for my own use and didn't account for other use-cases. If you need something more stable or a more mature library, check [r3labs/sse](https://github.com/r3labs/sse). Same goes for if you want a SSE Client library. But please feel free to tear the source and I'm open for PRs.

## Table of Contents

- [Synopsis](#synopsis)
- [Features](#feature)
- [Install](#install)
- [How to use it](#how-to-use-it)
- [Filtered/Selected Publish](#filteredselected-publish)
- [HTML/HTMX Templates](#htmlhtmx-templates)
- [Using Gin](#using-gin)

## Features

- Zero dependencies
- [melody](https://github.com/olahol/melody) inspired
- HTML/HTMX
- Compatible with standard mux handlers

## Install

```
go get github.com/christopher-kleine/sse@latest
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
		hub.ServeHTTP(w, r)
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
	return session.Get("role") == "villain"
})
```

## HTML/HTMX Templates

You can use this library to send HTML templates over SSE, since the `Event` type implements the `io.Writer` Interface:

```go
ev := &sse.Event{}
templates.ExecuteTemplate(ev, "index.html", nil)
hub.Publish(ev)
```

## Using Gin

The popular web framework [Gin](https://gin-gonic.com/) can be used too:

```go
package main

import (
	"net/http"
	"time"

	"github.com/christopher-kleine/sse"
	"github.com/gin-gonic/gin"
)

func main() {
	hub := sse.New()

	go func() {
		hub.Publish(&sse.Event{ Data: "ping" })
		time.Sleep(2 * time.Second)
	}()

	r := gin.Default()
	r.GET("/sse", func(c *gin.Context) {
		hub.ServeHTTP(c.Writer, c.Request)
	})
	r.Run()
}
```
