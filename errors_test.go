package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewWithCode(t *testing.T) {

	const (
		TESTERR ErrCode = iota
	)

	args := []interface{}{1, "2"}

	tests := []struct {
		code ErrCode
		err  string
		args []interface{}
		want error
	}{
		{NONE, "", nil, fmt.Errorf("")},
		{ERROR, "foo", nil, fmt.Errorf("foo")},
		{ERROR, "foo", nil, New("foo")},
		{ERROR, "foo", nil, NewWithCode(ERROR, "foo")},
		{TESTERR, "foo", nil, NewWithCode(ERROR, "foo")},
		{ERROR, "string with format specifiers: %v", nil, errors.New("string with format specifiers: %v")},
		{ERROR, "string with format specifiers: 1", nil, New("string with format specifiers: %v", 1)},
		{ERROR, "string with format specifiers: 1", nil, NewWithCode(ERROR, "string with format specifiers: %v", 1)},
		{TESTERR, "string with format specifiers: 1", nil, NewWithCode(ERROR, "string with format specifiers: %v", 1)},

		{NONE, "%d %s", args, fmt.Errorf("1 2")},
		{ERROR, "foo %d %s", args, fmt.Errorf("foo 1 2")},
		{ERROR, "foo %d %s", args, New("foo 1 2")},
		{ERROR, "foo %d %s", args, NewWithCode(ERROR, "foo 1 2")},
		{TESTERR, "foo %d %s", args, NewWithCode(ERROR, "foo 1 2")},
		{ERROR, "string with format specifiers: %d %s", args, errors.New("string with format specifiers: 1 2")},
		{ERROR, "string with format specifiers: %d %s", args, New("string with format specifiers: 1 2")},
		{ERROR, "string with format specifiers: %d %s", args, NewWithCode(ERROR, "string with format specifiers: 1 2")},
		{TESTERR, "string with format specifiers: %d %s", args, NewWithCode(ERROR, "string with format specifiers: 1 2")},
	}

	for _, tt := range tests {
		var got error
		if tt.args != nil {
			got = NewWithCode(tt.code, tt.err, tt.args...)
		} else {
			got = NewWithCode(tt.code, tt.err)
		}
		if got.Error() != tt.want.Error() {
			t.Errorf("NewWithCode.Error(): got: %q, want %q", got, tt.want)
		}
	}
}

func TestNew(t *testing.T) {

	args := []interface{}{1, "2"}

	tests := []struct {
		msg  string
		args []interface{}
		want error
	}{
		{"", nil, fmt.Errorf("")},
		{"foo", nil, fmt.Errorf("foo")},
		{"foo", nil, New("foo")},
		{"string with format specifiers: %v", nil, errors.New("string with format specifiers: %v")},
		{"%d %s", args, fmt.Errorf("1 2")},
		{"foo %d %s", args, fmt.Errorf("foo 1 2")},
		{"foo %d %s", args, New("foo 1 2")},
		{"string with format specifiers: %d %s", args, errors.New("string with format specifiers: 1 2")},
	}

	for _, tt := range tests {
		var got error
		if tt.args != nil {
			got = New(tt.msg, tt.args...)
		} else {
			got = New(tt.msg)
		}
		if got.Error() != tt.want.Error() {
			t.Errorf("New.Error(): got: %q, want %q", got, tt.want)
		}
	}
}

func TestAppendWithCode(t *testing.T) {

	const (
		TESTERR ErrCode = iota
	)

	args := []interface{}{1, "2"}

	var appendErr1 = errors.New("bar")
	var appendErr2 = NewWithCode(TESTERR, "bar")

	tests := []struct {
		code      ErrCode
		err       string
		args      []interface{}
		appendErr error
		want      error
	}{

		{ERROR, "", nil, appendErr1, fmt.Errorf("|bar")},
		{ERROR, "foo", nil, appendErr1, fmt.Errorf("foo|bar")},
		{ERROR, "foo", nil, appendErr1, New("foo|bar")},
		{ERROR, "foo", nil, appendErr1, NewWithCode(ERROR, "foo|bar")},
		{TESTERR, "foo", nil, appendErr1, NewWithCode(ERROR, "foo|bar")},
		{ERROR, "foo", nil, appendErr2, NewWithCode(ERROR, "foo|bar")},
		{ERROR, "string with format specifiers: %v", nil, appendErr1, errors.New("string with format specifiers: %v|bar")},
		{ERROR, "string with format specifiers: 1", nil, appendErr1, New("string with format specifiers: 1|bar")},
		{ERROR, "string with format specifiers: 1", nil, appendErr1, NewWithCode(ERROR, "string with format specifiers: 1|bar")},
		{TESTERR, "string with format specifiers: 1", nil, appendErr1, NewWithCode(ERROR, "string with format specifiers: 1|bar")},
		{ERROR, "string with format specifiers: 1", nil, appendErr2, NewWithCode(ERROR, "string with format specifiers: 1|bar")},

		{ERROR, "%d %s", args, appendErr1, fmt.Errorf("1 2|bar")},
		{ERROR, "foo %d %s", args, appendErr1, fmt.Errorf("foo 1 2|bar")},
		{ERROR, "foo %d %s", args, appendErr1, New("foo 1 2|bar")},
		{ERROR, "foo %d %s", args, appendErr1, NewWithCode(ERROR, "foo 1 2|bar")},
		{TESTERR, "foo %d %s", args, appendErr1, NewWithCode(ERROR, "foo 1 2|bar")},
		{ERROR, "foo %d %s", args, appendErr2, NewWithCode(ERROR, "foo 1 2|bar")},
		{ERROR, "string with format specifiers: %d %s", args, appendErr1, errors.New("string with format specifiers: 1 2|bar")},
		{ERROR, "string with format specifiers: %d %s", args, appendErr1, New("string with format specifiers: 1 2|bar")},
		{ERROR, "string with format specifiers: %d %s", args, appendErr1, NewWithCode(ERROR, "string with format specifiers: 1 2|bar")},
		{TESTERR, "string with format specifiers: %d %s", args, appendErr1, NewWithCode(ERROR, "string with format specifiers: 1 2|bar")},
		{ERROR, "string with format specifiers: %d %s", args, appendErr2, NewWithCode(ERROR, "string with format specifiers: 1 2|bar")},
	}

	for _, tt := range tests {
		var got error
		if tt.args != nil {
			got = AppendWithCode(tt.appendErr, tt.code, tt.err, tt.args...)
		} else {
			got = AppendWithCode(tt.appendErr, tt.code, tt.err)
		}

		if got.Error() != tt.want.Error() {
			t.Errorf("AppendWithCode.Error(): got: %q, want %q", got, tt.want)
		}
	}
}

func TestAppend(t *testing.T) {

	const (
		TESTERR ErrCode = iota
	)

	args := []interface{}{1, "2"}

	var appendErr1 = errors.New("bar")
	var appendErr2 = NewWithCode(TESTERR, "bar")

	tests := []struct {
		err       string
		args      []interface{}
		appendErr error
		want      error
	}{

		{"", nil, appendErr1, fmt.Errorf("|bar")},
		{"foo", nil, appendErr1, fmt.Errorf("foo|bar")},
		{"foo", nil, appendErr1, New("foo|bar")},
		{"foo", nil, appendErr1, NewWithCode(ERROR, "foo|bar")},
		{"foo", nil, appendErr2, NewWithCode(ERROR, "foo|bar")},
		{"string with format specifiers: %v", nil, appendErr1, errors.New("string with format specifiers: %v|bar")},
		{"string with format specifiers: 1", nil, appendErr1, New("string with format specifiers: 1|bar")},
		{"string with format specifiers: 1", nil, appendErr1, NewWithCode(ERROR, "string with format specifiers: 1|bar")},
		{"string with format specifiers: 1", nil, appendErr2, NewWithCode(ERROR, "string with format specifiers: 1|bar")},

		{"%d %s", args, appendErr1, fmt.Errorf("1 2|bar")},
		{"foo %d %s", args, appendErr1, fmt.Errorf("foo 1 2|bar")},
		{"foo %d %s", args, appendErr1, New("foo 1 2|bar")},
		{"foo %d %s", args, appendErr1, NewWithCode(ERROR, "foo 1 2|bar")},
		{"foo %d %s", args, appendErr2, NewWithCode(ERROR, "foo 1 2|bar")},
		{"string with format specifiers: %d %s", args, appendErr1, errors.New("string with format specifiers: 1 2|bar")},
		{"string with format specifiers: %d %s", args, appendErr1, New("string with format specifiers: 1 2|bar")},
		{"string with format specifiers: %d %s", args, appendErr1, NewWithCode(ERROR, "string with format specifiers: 1 2|bar")},
		{"string with format specifiers: %d %s", args, appendErr2, NewWithCode(ERROR, "string with format specifiers: 1 2|bar")},
	}

	for _, tt := range tests {
		var got error
		if tt.args != nil {
			got = Append(tt.appendErr, tt.err, tt.args...)
		} else {
			got = Append(tt.appendErr, tt.err)
		}

		if got.Error() != tt.want.Error() {
			t.Errorf("Append.Error(): got: %q, want %q", got, tt.want)
		}
	}
}

func TestNewWithCodeIs(t *testing.T) {

	const (
		TESTERR ErrCode = iota
	)

	args := []interface{}{1, "2"}

	tests := []struct {
		code ErrCode
		err  string
		args []interface{}
		want ErrCode
	}{
		{NONE, "", nil, NONE},
		{ERROR, "foo", nil, ERROR},
		{TESTERR, "foo", nil, TESTERR},
		{TESTERR, "string with format specifiers: %v", nil, TESTERR},

		{NONE, "%d %s", args, NONE},
		{ERROR, "foo %d %s", args, ERROR},
		{TESTERR, "foo %d %s", args, TESTERR},
	}

	for _, tt := range tests {
		var got error
		if tt.args != nil {
			got = NewWithCode(tt.code, tt.err, tt.args...)
		} else {
			got = NewWithCode(tt.code, tt.err)
		}
		if !Is(got, tt.want) {
			t.Errorf("New.Error(): got: %d, want %d", Code(got), tt.want)
		}
	}
}

func TestNewIs(t *testing.T) {

	args := []interface{}{1, "2"}

	tests := []struct {
		msg  string
		args []interface{}
		want ErrCode
	}{
		{"", nil, ERROR},
		{"foo", nil, ERROR},
		{"string with format specifiers: %v", nil, ERROR},
		{"%d %s", args, ERROR},
		{"foo %d %s", args, ERROR},
		{"string with format specifiers: %d %s", args, ERROR},
	}

	for _, tt := range tests {
		var got error
		if tt.args != nil {
			got = New(tt.msg, tt.args...)
		} else {
			got = New(tt.msg)
		}
		if !Is(got, tt.want) {
			t.Errorf("New.Error(): got: %d, want %d", Code(got), tt.want)
		}
	}
}

func TestAppendWithCodeIs(t *testing.T) {

	const (
		TESTERR ErrCode = iota
		TESTERRIN
	)

	args := []interface{}{1, "2"}

	var appendErr1 = errors.New("bar")
	var appendErr2 = NewWithCode(TESTERRIN, "bar")

	tests := []struct {
		code      ErrCode
		err       string
		args      []interface{}
		appendErr error
		want      ErrCode
	}{
		{TESTERR, "", nil, nil, TESTERR},
		{TESTERR, "foo", nil, nil, TESTERR},
		{TESTERR, "", nil, appendErr1, TESTERR},
		{TESTERR, "", nil, appendErr2, TESTERR},
		{TESTERR, "foo", nil, appendErr1, TESTERR},
		{TESTERR, "foo", nil, appendErr1, TESTERR},
		{TESTERR, "string with format specifiers: %v", nil, nil, TESTERR},

		{TESTERR, "%d %s", args, nil, TESTERR},
		{TESTERR, "foo %d %s", args, nil, TESTERR},
		{TESTERR, "%d %s", args, appendErr1, TESTERR},
		{TESTERR, "%d %s", args, appendErr2, TESTERR},
		{TESTERR, "foo %d %s", args, appendErr1, TESTERR},
		{TESTERR, "foo %d %s", args, appendErr2, TESTERR},
	}

	for _, tt := range tests {
		var got error
		if tt.args != nil {
			got = AppendWithCode(tt.appendErr, tt.code, tt.err, tt.args...)
		} else {
			got = AppendWithCode(tt.appendErr, tt.code, tt.err)
		}

		if !Is(got, tt.want) {
			t.Errorf("New.Error(): got: %d, want %d", Code(got), tt.want)
		}
	}
}

func TestAppendIs(t *testing.T) {

	const (
		TESTERR ErrCode = iota
	)

	args := []interface{}{1, "2"}

	var appendErr1 = errors.New("bar")
	var appendErr2 = NewWithCode(TESTERR, "bar")

	tests := []struct {
		err       string
		args      []interface{}
		appendErr error
		want      ErrCode
	}{
		{"", nil, nil, ERROR},
		{"foo", nil, nil, ERROR},
		{"", nil, appendErr1, ERROR},
		{"", nil, appendErr2, TESTERR},
		{"foo", nil, appendErr1, ERROR},
		{"foo", nil, appendErr2, TESTERR},
		{"string with format specifiers: %v", nil, nil, ERROR},
		{"string with format specifiers: %v", nil, appendErr1, ERROR},
		{"string with format specifiers: 1", nil, appendErr2, TESTERR},

		{"%d %s", args, nil, ERROR},
		{"foo %d %s", args, nil, ERROR},
		{"%d %s", args, appendErr1, ERROR},
		{"%d %s", args, appendErr2, TESTERR},
		{"foo %d %s", args, appendErr1, ERROR},
		{"foo %d %s", args, appendErr2, TESTERR},
	}

	for _, tt := range tests {
		var got error
		if tt.args != nil {
			got = Append(tt.appendErr, tt.err, tt.args...)
		} else {
			got = Append(tt.appendErr, tt.err, tt.args)
		}

		if !Is(got, tt.want) {
			t.Errorf("New.Error(): got: %d, want %d", Code(got), tt.want)
		}
	}
}
