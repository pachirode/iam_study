package validation

import (
	"fmt"
)

func MaxLenError(length int) string {
	return fmt.Sprintf("Must be no more than %d characters", length)
}

func RegexError(msg string, fmt string, examples ...string) string {
	if len(examples) == 0 {
		return msg + "(Regex used for validation is '" + fmt + "')"
	}

	msg += " (e.g. "
	for i := range examples {
		if i > 0 {
			msg += " or "
		}
		msg += "'" + examples[i] + "', "
	}
	msg += "regex used for validation is '" + fmt + "')"
	return msg
}

func InclusiveRangeError(lo, hi int) string {
	return fmt.Sprintf("Must between %d and %d", lo, hi)
}

func EmptyError() string {
	return "Must be non-empty"
}
