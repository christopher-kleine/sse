package main

import (
	sse "github.com/christopher-kleine/sse"
)

type App struct {
	hub *sse.Hub
}
