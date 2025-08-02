package field

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	utilerrors "github.com/pachirode/iam_study/pkg/errors"
	"github.com/pachirode/iam_study/pkg/sets"
)

type (
	ErrorType string
	ErrorList []*Error
)

const (
	ErrorTypeNotFound     ErrorType = "FieldValueNotFound"
	ErrorTypeRequired     ErrorType = "FieldValueRequired"
	ErrorTypeDuplicate    ErrorType = "FieldValueDuplicate"
	ErrorTypeInvalid      ErrorType = "FieldValueInvalid"
	ErrorTypeNotSupported ErrorType = "FieldValueNotSupported"
	ErrorTypeForbidden    ErrorType = "FieldValueForbidden"
	ErrorTypeTooLong      ErrorType = "FieldValueTooLong"
	ErrorTypeTooMany      ErrorType = "FieldValueTooMany"
	ErrorTypeInternal     ErrorType = "InternalError"
)

type Error struct {
	Type     ErrorType
	Field    string
	BadValue interface{}
	Detail   string
}

func (err *Error) Error() string {
	return fmt.Sprintf("%s: %s", err.Field, err.ErrorBody())
}

func (err *Error) ErrorBody() string {
	var s string
	switch err.Type {
	case ErrorTypeRequired, ErrorTypeForbidden, ErrorTypeTooLong, ErrorTypeInternal:
		s = err.Type.String()
	default:
		value := err.BadValue
		valueType := reflect.TypeOf(value)
		if value == nil || valueType == nil {
			value = "null"
		} else if valueType.Kind() == reflect.Ptr {
			if reflectValue := reflect.ValueOf(value); reflectValue.IsNil() {
				value = "null"
			} else {
				value = reflectValue.Elem().Interface()
			}
		}
		switch t := value.(type) {
		case int64, int32, float64, float32, bool:
			s = fmt.Sprintf("%s: %v", err.Type, value)
		case string:
			s = fmt.Sprintf("%s: %q", err.Type, t)
		case fmt.Stringer:
			s = fmt.Sprintf("%s: %s", err.Type, t.String())
		default:
			s = fmt.Sprintf("%s: %#v", err.Type, value)
		}
	}

	if len(err.Detail) > 0 {
		s += fmt.Sprintf(": %s", err.Detail)
	}

	return s
}

func (t ErrorType) String() string {
	switch t {
	case ErrorTypeNotFound:
		return "Not found"
	case ErrorTypeRequired:
		return "Required value"
	case ErrorTypeDuplicate:
		return "Duplicate value"
	case ErrorTypeInvalid:
		return "Invalid value"
	case ErrorTypeNotSupported:
		return "Unsupported value"
	case ErrorTypeForbidden:
		return "Forbidden"
	case ErrorTypeTooLong:
		return "Too long"
	case ErrorTypeTooMany:
		return "Too many"
	case ErrorTypeInternal:
		return "Internal error"
	default:
		panic(fmt.Sprintf("unrecognized validation error: %q", string(t)))
	}
}

func NotFound(field *Path, value interface{}) *Error {
	return &Error{ErrorTypeNotFound, field.String(), value, ""}
}

func Required(field *Path, detail string) *Error {
	return &Error{ErrorTypeRequired, field.String(), "", detail}
}

func Duplicate(field *Path, value interface{}) *Error {
	return &Error{ErrorTypeDuplicate, field.String(), value, ""}
}

func Invalid(field *Path, value interface{}, detail string) *Error {
	return &Error{ErrorTypeInvalid, field.String(), value, detail}
}

func NotSupported(field *Path, value interface{}, validValues []string) *Error {
	detail := ""
	if len(validValues) > 0 {
		quotedValues := make([]string, len(validValues))
		for i, v := range validValues {
			quotedValues[i] = strconv.Quote(v)
		}
		detail = "supported values: " + strings.Join(quotedValues, ", ")
	}
	return &Error{ErrorTypeNotSupported, field.String(), value, detail}
}

func Forbidden(field *Path, detail string) *Error {
	return &Error{ErrorTypeForbidden, field.String(), "", detail}
}

func TooLong(field *Path, value interface{}, maxLength int) *Error {
	return &Error{ErrorTypeTooLong, field.String(), value, fmt.Sprintf("must have at most %d bytes", maxLength)}
}

func TooMany(field *Path, actualQuantity, maxQuantity int) *Error {
	return &Error{
		ErrorTypeTooMany,
		field.String(),
		actualQuantity,
		fmt.Sprintf("must have at most %d items", maxQuantity),
	}
}

func InternalError(field *Path, err error) *Error {
	return &Error{ErrorTypeInternal, field.String(), nil, err.Error()}
}

func NewErrorTypeMatcher(t ErrorType) utilerrors.Matcher {
	return func(err error) bool {
		if e, ok := err.(*Error); ok {
			return e.Type == t
		}
		return false
	}
}

func (list ErrorList) ToAggregate() utilerrors.Aggregate {
	errs := make([]error, 0, len(list))
	errorMsgs := sets.NewString()
	for _, err := range list {
		msg := fmt.Sprintf("%v", err)
		if errorMsgs.Has(msg) {
			continue
		}
		errorMsgs.Insert(msg)
		errs = append(errs, err)
	}
	return utilerrors.NewAggregate(errs)
}

func fromAggregate(agg utilerrors.Aggregate) ErrorList {
	errs := agg.Errors()
	list := make(ErrorList, len(errs))
	for i := range errs {
		list[i] = errs[i].(*Error)
	}
	return list
}

func (list ErrorList) Filter(fns ...utilerrors.Matcher) ErrorList {
	err := utilerrors.FilterOut(list.ToAggregate(), fns...)
	if err == nil {
		return nil
	}
	// FilterOut takes an Aggregate and returns an Aggregate
	return fromAggregate(err.(utilerrors.Aggregate))
}
