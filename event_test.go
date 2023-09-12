package sse_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/christopher-kleine/sse"
)

func TestEventString(t *testing.T) {
	testTable := []struct {
		name     string
		input    *sse.Event
		expected string
	}{
		{name: "data-only", input: &sse.Event{Data: "test"}, expected: "data: test\n\n"},
		{name: "retry-only", input: &sse.Event{Retry: 1 * time.Second}, expected: "retry: 1000\n\n"},
		{name: "id-only", input: &sse.Event{ID: "1"}, expected: "id: 1\n\n"},
		{name: "event-only", input: &sse.Event{Event: "foo"}, expected: "event: foo\n\n"},
		{name: "rety-and-data", input: &sse.Event{Retry: 1 * time.Second, Data: "Dummy"}, expected: "data: Dummy\nretry: 1000\n\n"},
		{name: "multiline-data", input: &sse.Event{Data: "foo\nbar"}, expected: "data: foo\ndata: bar\n\n"},
		{name: "multiline-data-with-eol", input: &sse.Event{Data: "foo\nbar\n"}, expected: "data: foo\ndata: bar\n\n"},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.input.String()
			if actual != testCase.expected {
				t.Errorf("%q should be %q", actual, testCase.expected)
			}
		})
	}
}

func TestEventWrite(t *testing.T) {
	testTable := []struct {
		name     string
		input    io.Reader
		expected string
	}{
		{
			name:     "simple-text",
			input:    bytes.NewBufferString("foo"),
			expected: "data: foo\n\n",
		},
		{
			name:     "multiline-text",
			input:    bytes.NewBufferString("foo\nbar"),
			expected: "data: foo\ndata: bar\n\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ev := &sse.Event{}
			io.Copy(ev, testCase.input)
			actual := ev.String()
			if actual != testCase.expected {
				t.Errorf("%q should be %q", actual, testCase.expected)
			}
		})
	}
}
