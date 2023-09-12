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

func TestEventJSONData(t *testing.T) {
	testTable := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "string",
			input:    "foo",
			expected: "data: \"foo\"\n\n",
		},
		{
			name:     "integer",
			input:    1,
			expected: "data: 1\n\n",
		},
		{
			name:     "float",
			input:    3.14159,
			expected: "data: 3.14159\n\n",
		},
		{
			name:     "bool (true)",
			input:    true,
			expected: "data: true\n\n",
		},
		{
			name:     "bool (false)",
			input:    false,
			expected: "data: false\n\n",
		},
		{
			name:     "nil",
			input:    nil,
			expected: "data: null\n\n",
		},
		{
			name:     "struct-1",
			input:    struct{ Text string }{Text: "foo"},
			expected: "data: {\"Text\":\"foo\"}\n\n",
		},
		{
			name: "struct-2",
			input: struct {
				Text string `json:"text"`
			}{Text: "foo"},
			expected: "data: {\"text\":\"foo\"}\n\n",
		},
		{
			name: "struct-3",
			input: struct {
				Text string  `json:"text"`
				Bar  *string `json:"bar,omitempty"`
			}{Text: "foo"},
			expected: "data: {\"text\":\"foo\"}\n\n",
		},
		{
			name: "struct-4",
			input: struct {
				Text string  `json:"text"`
				Bar  *string `json:"bar"`
			}{Text: "foo"},
			expected: "data: {\"text\":\"foo\",\"bar\":null}\n\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ev := &sse.Event{}
			ev.JSONData(testCase.input)
			actual := ev.String()
			if actual != testCase.expected {
				t.Errorf("%q should be %q", actual, testCase.expected)
			}
		})
	}
}
