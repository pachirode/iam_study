package errors

import (
	"fmt"
	"net/http"
	"sync"
)

type Coder interface {
	HTTPStatus() int
	String() string
	Reference() string
	Code() int
}

type coder struct {
	HTTPNum int
	Msg     string
	Ref     string
	CodeNum int
}

var (
	_            Coder = (*coder)(nil)
	codes              = map[int]Coder{}
	codeMux            = &sync.Mutex{}
	unknownCoder       = coder{http.StatusInternalServerError, "Internal server error", "", 1}
)

func (c coder) Code() int {
	return c.CodeNum
}

func (c coder) String() string {
	return c.Msg
}

func (c coder) Reference() string {
	return c.Ref
}

func (c coder) HTTPStatus() int {
	if c.HTTPNum == 0 {
		return 500
	}

	return c.HTTPNum
}

func Register(c coder) {
	if c.Code() == 0 {
		panic("code `0` is used by default")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	codes[c.Code()] = c
}

func MustRegister(c coder) {
	if c.Code() == 0 {
		panic("code `0` is used by default")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[c.Code()]; ok {
		panic(fmt.Sprintf("code: %d already exist", c.Code()))
	}

	codes[c.Code()] = c
}

type withCode struct {
	err   error
	code  int
	cause error
	*stack
}

func (w *withCode) Error() string {
	return fmt.Sprintf("%v", w)
}

func (w *withCode) Cause() error {
	return w.cause
}

func (w *withCode) Unwrap() error {
	return w.cause
}

func WithCode(code int, format string, args ...interface{}) error {
	return &withCode{
		err:   fmt.Errorf(format, args...),
		code:  code,
		stack: callers(),
	}
}

func WrapCode(err error, code int, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &withCode{
		err:   fmt.Errorf(format, args...),
		code:  code,
		cause: err,
		stack: callers(),
	}
}
