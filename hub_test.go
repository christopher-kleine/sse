package sse_test

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/christopher-kleine/sse"
)

func TestHubPublish(t *testing.T) {
	expected := `data: foo
id: 0

data: foo
id: 1

data: foo
id: 2

data: foo
id: 3

data: foo
id: 4

`

	// Create a new hub and a Testserver
	hub := sse.New()
	ts := httptest.NewServer(hub)
	defer ts.Close()

	// Create a new Request and add a Cancel function
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(req.Context())
	req = req.WithContext(ctx)

	// Wait for 10 Milliseconds between each publish
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(10 * time.Millisecond)
			hub.Publish(&sse.Event{Data: "foo", ID: strconv.Itoa(i)})
		}
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	// Connect to the Testserver and prepare the variables
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	reader := bufio.NewReader(res.Body)
	buffer := make([]byte, 4096)
	var testResult string

	// Read as long as needed
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			break
		}
		testResult += string(buffer[:n])
	}

	if testResult != expected {
		t.Errorf("%q should be %q", testResult, expected)
	}
}
