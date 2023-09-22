package sse_test

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/christopher-kleine/sse"
)

func TestHubPublish(t *testing.T) {
	hub := sse.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx, cancel := context.WithCancel(req.Context())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	go func() {
		reader := bufio.NewReader(w.Body)
		hub.ServeHTTP(w, req)

		var err error
		var line string
		for err != io.EOF {
			line, err = reader.ReadString('\n')
			t.Logf("Date: %s, Line: %q, err: %v", time.Now(), line, err)
		}
	}()

	time.Sleep(1 * time.Second)
	hub.Publish(&sse.Event{Data: "Test"})

	time.Sleep(5 * time.Second)
	hub.Publish(&sse.Event{Data: "Foo"})

	time.Sleep(5 * time.Second)
	cancel()

	time.Sleep(1 * time.Second)
}
