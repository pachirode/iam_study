package errors

import "errors"

type MessagaCountMap map[string]int

type Aggregate interface {
	error
	Errors() []error
	Is(error) bool
}

type aggregate []error

type Matcher func(error) bool

func (agg aggregate) Is(target error) bool {
	return agg.visit(func(err error) bool {
		return errors.Is(err, target)
	})
}

func (agg aggregate) visit(f func(err error) bool) bool {
	for _, err := range agg {
		switch err := err.(type) {
		case aggregate:
			if match := err.visit(f); match {
				return match
			}
		case Aggregate:
			for _, nestedErr := range err.Errors() {
				if match := f(nestedErr); match {
					return match
				}
			}
		default:
			if match := f(err); match {
				return match
			}
		}
	}

	return false
}

func (agg aggregate) Errors() []error {
	return agg
}

func NewAggregate(errlist []error) Aggregate {
	if len(errlist) == 0 {
		return nil
	}

	var errs []error
	for _, e := range errlist {
		if e != nil {
			errs = append(errs, e)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return nil
}

func (agg aggregate) Error() string {
	if len(agg) == 0 {
		return ""
	}

	if len(agg) == 1 {
		return agg[0].Error()
	}

	seenerrs := NewString()
	res := ""

	agg.visit(func(err error) bool {
		msg := err.Error()
		if seenerrs.Has(msg) {
			return false
		}
		seenerrs.Insert(msg)
		if len(seenerrs) > 1 {
			res += ", "
		}
		res += msg
		return true
	})
	if len(seenerrs) == 1 {
		return res
	}

	return "[" + res + "]"
}

func FilterOut(err error, fns ...Matcher) error {
	if err == nil {
		return nil
	}

	if agg, ok := err.(Aggregate); ok {
		return NewAggregate(filterErrors(agg.Errors(), fns...))
	}

	if !matchesError(err, fns...) {
		return err
	}

	return nil
}

func matchesError(err error, fns ...Matcher) bool {
	for _, fn := range fns {
		if fn(err) {
			return true
		}
	}

	return false
}

func filterErrors(list []error, fns ...Matcher) []error {
	res := []error{}
	for _, err := range list {
		r := FilterOut(err, fns...)
		if r != nil {
			res = append(res, r)
		}
	}
	return res
}
