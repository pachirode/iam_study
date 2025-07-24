package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strings"
)

type Frame uintptr

func (f Frame) pc() uintptr {
	return uintptr(f) - 1
}

func (f Frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unkonw"
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

func (f Frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

func (f Frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unkonw"
	}
	return fn.Name()
}

func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			fmt.Fprintf(s, "%s\n\t%s", f.name(), f.file())
		default:
			fmt.Fprintf(s, "%s", path.Base(f.file()))
		}
	case 'd':
		fmt.Fprintf(s, "%d", f.line())
	case 'n':
		fmt.Fprintf(s, "%s", funcname(f.name()))
	case 'v':
		fmt.Fprintf(s, "%s:%d", f, f.line())
	}
}

func (f Frame) MarshalText() ([]byte, error) {
	name := f.name()
	if name == "unkonw" {
		return []byte(name), nil
	}
	return []byte(fmt.Sprintf("%s %s:%d", name, f.file(), f.line())), nil
}

type StackTrace []Frame

func (st StackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			for _, f := range st {
				io.WriteString(s, "\n")
				f.Format(s, verb)
			}
		case s.Flag('#'):
			fmt.Fprintf(s, "%#v", []Frame(st))
		default:
			st.formatSlice(s, verb)
		}
	case 's':
		st.formatSlice(s, verb)
	}
}

func (st StackTrace) formatSlice(s fmt.State, verb rune) {
	io.WriteString(s, "[")
	for i, f := range st {
		if i > 0 {
			io.WriteString(s, " ")
		}
		f.Format(s, verb)
	}
	io.WriteString(s, "]")
}

type stack []uintptr

func (st *stack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			for _, pc := range *st {
				f := Frame(pc)
				fmt.Fprintf(s, "\n%+v", f)
			}
		}
	}
}

func (st *stack) StackTrace() StackTrace {
	f := make([]Frame, len(*st))
	for i := 0; i < len(f); i++ {
		f[i] = Frame((*st)[i])
	}
	return f
}

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}

type withStack struct {
	error
	*stack
}

func (w *withStack) Cause() error {
	return w.error
}

func (w *withStack) Unwrap() error {
	if e, ok := w.error.(interface{ Unwrap() error }); ok {
		return e.Unwrap()
	}

	return w.error
}

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}
