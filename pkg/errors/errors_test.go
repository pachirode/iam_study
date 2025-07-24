package errors

import (
	"fmt"
	"io"
	"reflect"
	"testing"
)

type TestInfo struct {
	err  string
	want error
}

type TestWrapInfo struct {
	err     error
	message string
	want    string
}

type TestError struct {
	err  error
	want error
}

func TestNew(t *testing.T) {
	tests := []TestInfo{
		{"", fmt.Errorf("")},
		{"foo", fmt.Errorf("foo")},
		{"foo", New("foo")},
	}

	for _, test := range tests {
		got := New(test.err)
		if got.Error() != test.want.Error() {
			t.Errorf("New.Error(): got: %q, want: %q", got, test.want)
		}
	}
}

func TestWrapNil(t *testing.T) {
	got := Wrap(nil, "msg")
	if got != nil {
		t.Errorf("Wrap(nil, \"no error\"): got %#v, excepted nil", got)
	}
}

func TestWrap(t *testing.T) {
	tests := []TestWrapInfo{
		{io.EOF, "read error", "read error"},
		{Wrap(io.EOF, "read error"), "client error", "client error"},
	}

	for _, test := range tests {
		got := Wrap(test.err, test.message).Error()
		if got != test.want {
			t.Errorf("Wrap(%v, %q): got: %v, want: %v", test.err, test.message, got, test.want)
		}
	}
}

func TestCause(t *testing.T) {
	tests := []TestError{
		{
			err:  nil,
			want: nil,
		},
		{
			err:  (error)(nil),
			want: nil,
		},
		{
			err:  io.EOF,
			want: io.EOF,
		},
	}

	for i, test := range tests {
		got := Cause(test.err)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("test %d: got %#v, want %#v", i+1, got, test.want)
		}
	}
}

func TestWrapf(t *testing.T) {
	tests := []TestWrapInfo{
		{io.EOF, "read error", "read error"},
		{Wrap(io.EOF, "read error"), "client error", "client error"},
	}
	for _, test := range tests {
		got := Wrapf(test.err, test.message).Error()
		if got != test.want {
			t.Errorf("Wrapf(%v, %q): got: %v, want: %v", test.err, test.message, got, test.want)
		}
	}
}

func TestErrorf(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{Errof("read error without format specifiers"), "read error without format specifiers"},
		{Errof("read error with %d format", 1), "read error with 1 format"},
	}

	for _, test := range tests {
		got := test.err.Error()
		if got != test.want {
			t.Errorf("Errorf(%v), got: %q, want: %q", test.err, got, test.want)
		}
	}
}
