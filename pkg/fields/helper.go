package fields

import (
	"strings"
)

const (
	notEqualOperator    = "!="
	doubleEqualOperator = "=="
	equalOperator       = "="
)

var (
	termOperators = []string{notEqualOperator, doubleEqualOperator, equalOperator}
	valueEscaper  = strings.NewReplacer(
		`\`, `\\`,
		`,`, `\,`,
		`=`, `\=`,
	)
)

func splitTerms(fieldSelector string) []string {
	if len(fieldSelector) == 0 {
		return nil
	}

	terms := make([]string, 0, 1)
	startIndex := 0
	isSlash := false
	for i, c := range fieldSelector {
		switch {
		case isSlash:
			isSlash = false
		case c == '\\':
			isSlash = true
		case c == ',':
			terms = append(terms, fieldSelector[startIndex:i])
			startIndex = i + 1
		}
	}

	terms = append(terms, fieldSelector[startIndex:])

	return terms
}

func splitTerm(term string) (lhs, op, rhs string, ok bool) {
	for i := range term {
		remaining := term[i:]
		for _, op := range termOperators {
			if strings.HasPrefix(remaining, op) {
				return term[0:i], op, term[i+len(op):], true
			}
		}
	}

	return "", "", "", false
}

func EscapeValue(s string) string {
	return valueEscaper.Replace(s)
}
