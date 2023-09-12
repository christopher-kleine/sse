package sse_test

import (
	"testing"
	"time"

	"github.com/christopher-kleine/sse"
)

func TestEventString(t *testing.T) {
	testTable := []struct {
		input    *sse.Event
		expected string
	}{
		{input: &sse.Event{Data: "test"}, expected: "data: test\n\n"},
		{input: &sse.Event{Retry: 1 * time.Second}, expected: "retry: 1000\n\n"},
		{input: &sse.Event{ID: "1"}, expected: "id: 1\n\n"},
		{input: &sse.Event{Event: "foo"}, expected: "event: foo\n\n"},
		{input: &sse.Event{Retry: 1 * time.Second, Data: "Dummy"}, expected: "data: Dummy\nretry: 1000\n\n"},
		{input: &sse.Event{Data: "foo\nbar"}, expected: "data: foo\ndata: bar\n\n"},
		{input: &sse.Event{Data: "foo\nbar\n"}, expected: "data: foo\ndata: bar\n\n"},
	}

	for _, testCase := range testTable {
		actual := testCase.input.String()
		if actual != testCase.expected {
			t.Errorf("%q should be %q", actual, testCase.expected)
		}
	}
}
