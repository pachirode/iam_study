package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type formatInfo struct {
	code    int
	message string
	err     string
	stack   *stack
}

func (w *withCode) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		str := bytes.NewBuffer([]byte{})
		jsonData := []map[string]interface{}{}

		var (
			flagDetail bool
			flagTrace  bool
			modeJSON   bool
		)

		if state.Flag('#') {
			modeJSON = true
		}

		if state.Flag('-') {
			flagDetail = true
		}
		if state.Flag('+') {
			flagTrace = true
		}

		sep := ""
		errs := list(w)
		length := len(errs)
		for k, e := range errs {
			finfo := buildFormatInfo(e)
			jsonData, str = format(length-k-1, jsonData, str, finfo, sep, flagDetail, flagTrace, modeJSON)
			sep = "; "

			if !flagTrace {
				break
			}

			if !flagDetail && !flagTrace && !modeJSON {
				break
			}
		}
		if modeJSON {
			var byts []byte
			byts, _ = json.Marshal(jsonData)

			str.Write(byts)
		}

		fmt.Fprintf(state, "%s", strings.Trim(str.String(), "\r\n\t"))
	default:
		finfo := buildFormatInfo(w)
		fmt.Fprintf(state, finfo.message)
	}
}

func format(k int, jsonData []map[string]interface{}, str *bytes.Buffer, finfo *formatInfo,
	sep string, flagDetail, flagTrace, modeJSON bool,
) ([]map[string]interface{}, *bytes.Buffer) {
	if modeJSON {
		data := map[string]interface{}{}
		if flagDetail || flagTrace {
			data = map[string]interface{}{
				"message": finfo.message,
				"code":    finfo.code,
				"error":   finfo.err,
			}

			caller := fmt.Sprintf("#%d", k)
			if finfo.stack != nil {
				f := Frame((*finfo.stack)[0])
				caller = fmt.Sprintf("%s %s:%d (%s)",
					caller,
					f.file(),
					f.line(),
					f.name(),
				)
			}
			data["caller"] = caller
		} else {
			data["error"] = finfo.message
		}
		jsonData = append(jsonData, data)
	} else {
		if flagDetail || flagTrace {
			if finfo.stack != nil {
				f := Frame((*finfo.stack)[0])
				fmt.Fprintf(str, "%s%s - #%d [%s:%d (%s)] (%d) %s",
					sep,
					finfo.err,
					k,
					f.file(),
					f.line(),
					f.name(),
					finfo.code,
					finfo.message,
				)
			} else {
				fmt.Fprintf(str, "%s%s - #%d %s", sep, finfo.err, k, finfo.message)
			}
		} else {
			fmt.Fprintf(str, finfo.message)
		}
	}

	return jsonData, str
}

func list(e error) []error {
	ret := []error{}

	if e != nil {
		if w, ok := e.(interface{ Unwrap() error }); ok {
			ret = append(ret, e)
			ret = append(ret, list(w.Unwrap())...)
		} else {
			ret = append(ret, e)
		}
	}

	return ret
}

func buildFormatInfo(e error) *formatInfo {
	var finfo *formatInfo

	switch err := e.(type) {
	case *fundamental:
		finfo = &formatInfo{
			code:    unknownCoder.Code(),
			message: err.msg,
			err:     err.msg,
			stack:   err.stack,
		}
	case *withStack:
		finfo = &formatInfo{
			code:    unknownCoder.Code(),
			message: err.Error(),
			err:     err.Error(),
			stack:   err.stack,
		}
	case *withCode:
		coder, ok := codes[err.code]
		if !ok {
			coder = unknownCoder
		}

		msg := coder.String()
		if msg == "" {
			msg = err.err.Error()
		}

		finfo = &formatInfo{
			code:    coder.Code(),
			message: msg,
			err:     err.err.Error(),
			stack:   err.stack,
		}
	default:
		finfo = &formatInfo{
			code:    unknownCoder.Code(),
			message: err.Error(),
			err:     err.Error(),
		}
	}

	return finfo
}
