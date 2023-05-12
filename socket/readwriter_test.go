package socket

import (
	"bytes"
	"sync"
	"testing"
)

func TestRead(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       string
		expectedEv  Event
		expectedArg string
	}{
		{
			name:        "Test ADD",
			input:       "ADD player1",
			expectedEv:  Add,
			expectedArg: "player1",
		},
		{
			name:        "Test MATCH",
			input:       "MATCH some_match123",
			expectedEv:  Match,
			expectedArg: "some_match123",
		},
		{
			name:        "Test REMOVE",
			input:       "REMOVE some_player",
			expectedEv:  Remove,
			expectedArg: "some_player",
		},
		{
			name:        "Test SIZE",
			input:       "SIZE",
			expectedEv:  Size,
			expectedArg: "",
		},
		{
			name:        "Test SIZE",
			input:       "SIZE 6",
			expectedEv:  Size,
			expectedArg: "6",
		},
		{
			name:        "Test Invalid Event",
			input:       "INVALID_EVENT some_arg",
			expectedEv:  Unknown,
			expectedArg: "INVALID_EVENT some_arg",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			buf.WriteString(tc.input + "\n")
			rw := newReadWriter(&buf)

			event, args, err := rw.Read()
			if err != nil {
				t.Errorf("error read: %v", err)
			}

			if event != tc.expectedEv {
				t.Errorf("event: want %s got %s", tc.expectedEv, event)
			}

			if args != tc.expectedArg {
				t.Errorf("args: want %s got %s", tc.expectedArg, args)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		event    Event
		args     []interface{}
		expected string
		err      error
	}{
		{
			name:     "event without args",
			event:    Add,
			args:     nil,
			expected: "ADD\n",
			err:      nil,
		},
		{
			name:     "event with single arg",
			event:    Remove,
			args:     []interface{}{"foo"},
			expected: "REMOVE foo\n",
			err:      nil,
		},
		{
			name:     "event with multiple args",
			event:    Match,
			args:     []interface{}{123, "bar"},
			expected: "MATCH 123 bar\n",
			err:      nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var b bytes.Buffer
			rw := newReadWriter(&b)

			if err := rw.Write(tc.event, tc.args...); err != nil {
				t.Errorf("error write: %v", err)
			}

			s, err := b.ReadString(Delimiter)
			if err != nil {
				t.Errorf("error read: %v", err)
			}

			if s != tc.expected {
				t.Errorf("result: want %s got %s", tc.expected, s)
			}
		})
	}
}

func TestWriteConcurrent(t *testing.T) {
	t.Parallel()
	var (
		b  bytes.Buffer
		rw = newReadWriter(&b)
		i  = 100_000
		wg sync.WaitGroup
	)

	wg.Add(i)

	for j := 0; j < i; j++ {
		go func() {
			if err := rw.Write(Unknown); err != nil {
				t.Errorf("error write: %v", err)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	for j := 0; j < i; j++ {
		if _, err := b.ReadBytes(Delimiter); err != nil {
			t.Errorf("error read: %v", err)
		}
	}
}
