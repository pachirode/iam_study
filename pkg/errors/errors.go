package errors

import (
	"fmt"
	"io"
)

func New(message string) error {
	return &fundamental{
		msg:   message,
		stack: callers(),
	}
}

func Errof(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

type fundamental struct {
	msg string
	*stack
}

func (f *fundamental) Error() string {
	return f.msg
}

func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, f.msg)
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, f.msg)
	case 'q':
		fmt.Fprintf(s, "%q", f.msg)
	}
}

func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}

	return &withMessage{
		cause: err,
		msg:   message,
	}
}

func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string {
	return w.msg
}

func (w *withMessage) Cause() error {
	return w.cause
}

// for Go 1.13
func (w *withMessage) Unwrap() error {
	return w.cause
}

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*withCode); ok {
		return &withCode{
			err:   fmt.Errorf(message),
			code:  e.code,
			cause: err,
			stack: callers(),
		}
	}

	err = &withMessage{
		cause: err,
		msg:   message,
	}

	return &withStack{
		err,
		callers(),
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*withCode); ok {
		return &withCode{
			err:   fmt.Errorf(format, args...),
			code:  e.code,
			cause: err,
			stack: callers(),
		}
	}

	err = &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}

	return &withStack{
		err,
		callers(),
	}
}

func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}

		if cause.Cause() == nil {
			break
		}

		err = cause.Cause()
	}

	return err
}
