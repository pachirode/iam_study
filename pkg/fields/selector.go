package fields

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/pachirode/iam_study/pkg/selection"
)

type Selector interface {
	Matches(Fields) bool
	Empty() bool
	RequiresExactMatch(field string) (value string, found bool)
	Transform(fn TransformFunc) (Selector, error)
	Requirements() Requirements
	String() string
	DeepCopySelector() Selector
}

type nothingSelector struct{}

type hasTerm struct {
	field, value string
}

type notHasTerm struct {
	field, value string
}

type InvalidEscapeSequence struct {
	sequence string
}

type UnescapedRune struct {
	r rune
}

type andTerm []Selector

type TransformFunc func(field, value string) (newField, newValue string, err error)

func (n nothingSelector) Matches(_ Fields) bool {
	return false
}

func (n nothingSelector) Empty() bool {
	return false
}

func (n nothingSelector) String() string {
	return ""
}

func (n nothingSelector) Requirements() Requirements {
	return nil
}

func (n nothingSelector) DeepCopySelector() Selector {
	return n
}

func (n nothingSelector) RequiresExactMatch(field string) (value string, found bool) {
	return "", false
}

func (n nothingSelector) Transform(fn TransformFunc) (Selector, error) {
	return n, nil
}

func (t *hasTerm) Matches(ls Fields) bool {
	return ls.Get(t.field) == t.value
}

func (t *hasTerm) Empty() bool {
	return false
}

func (t *hasTerm) RequiresExactMatch(field string) (value string, found bool) {
	if t.field == field {
		return t.value, true
	}

	return "", false
}

func (t *hasTerm) Transform(fn TransformFunc) (Selector, error) {
	field, value, err := fn(t.field, t.value)
	if err != nil {
		return nil, err
	}

	if len(field) == 0 && len(value) == 0 {
		return Everything(), nil
	}

	return &hasTerm{field: field, value: value}, nil
}

func (t *hasTerm) Requirements() Requirements {
	return []Requirement{{
		Field:    t.field,
		Operator: selection.Equals,
		Value:    t.value,
	}}
}

func (t *hasTerm) String() string {
	return fmt.Sprintf("%v=%v", t.field, EscapeValue(t.value))
}

func (t *hasTerm) DeepCopySelector() Selector {
	if t == nil {
		return nil
	}

	out := new(hasTerm)
	*out = *t

	return out
}

func (t *notHasTerm) Matches(ls Fields) bool {
	return ls.Get(t.field) != t.value
}

func (t *notHasTerm) Empty() bool {
	return false
}

func (t *notHasTerm) RequiresExactMatch(field string) (value string, found bool) {
	return "", false
}

func (t *notHasTerm) Transform(fn TransformFunc) (Selector, error) {
	field, value, err := fn(t.field, t.value)
	if err != nil {
		return nil, err
	}

	if len(field) == 0 && len(value) == 0 {
		return Everything(), nil
	}

	return &notHasTerm{field: field, value: value}, nil
}

func (t *notHasTerm) Requirements() Requirements {
	return []Requirement{{
		Field:    t.field,
		Operator: selection.NoEquals,
		Value:    t.value,
	}}
}

func (t *notHasTerm) String() string {
	return fmt.Sprintf("%v=%v", t.field, EscapeValue(t.value))
}

func (t *notHasTerm) DeepCopySelector() Selector {
	if t == nil {
		return nil
	}

	out := new(notHasTerm)
	*out = *t

	return out
}

func (t andTerm) Matches(ls Fields) bool {
	for _, q := range t {
		if !q.Matches(ls) {
			return false
		}
	}

	return true
}

func (t andTerm) Empty() bool {
	if t == nil {
		return true
	}
	if len([]Selector(t)) == 0 {
		return true
	}

	for i := range t {
		if !t[i].Empty() {
			return false
		}
	}

	return true
}

func (t andTerm) RequiresExactMatch(field string) (string, bool) {
	if t == nil || len([]Selector(t)) == 0 {
		return "", false
	}

	for i := range t {
		if value, found := t[i].RequiresExactMatch(field); found {
			return value, found
		}
	}

	return "", false
}

func (t andTerm) Transform(fn TransformFunc) (Selector, error) {
	next := make([]Selector, 0, len([]Selector(t)))
	for _, s := range []Selector(t) {
		n, err := s.Transform(fn)
		if err != nil {
			return nil, err
		}
		if !n.Empty() {
			next = append(next, n)
		}
	}

	return andTerm(next), nil
}

func (t andTerm) Requirements() Requirements {
	reqs := make([]Requirement, 0, len(t))
	for _, s := range []Selector(t) {
		rs := s.Requirements()
		reqs = append(reqs, rs...)
	}

	return reqs
}

func (t andTerm) String() string {
	terms := make([]string, 0, len(t))
	for _, q := range t {
		terms = append(terms, q.String())
	}

	return strings.Join(terms, ",")
}

func (t andTerm) DeepCopySelector() Selector {
	if t == nil {
		return nil
	}
	out := make([]Selector, 0, len(t))
	for i := range t {
		out[i] = t[i].DeepCopySelector()
	}

	return andTerm(out)
}

func (i InvalidEscapeSequence) Error() string {
	return fmt.Sprintf("Invalid field selector: invalid escape sequence: %s", i.sequence)
}

func (u UnescapedRune) Error() string {
	return fmt.Sprintf("Invalid field selector: unescaped character in value: %v", u.r)
}

func SelectorFromSet(ls Set) Selector {
	if ls == nil {
		return Everything()
	}

	items := make([]Selector, 0, len(ls))

	for field, value := range ls {
		items = append(items, &hasTerm{field: field, value: value})
	}

	if len(items) == 1 {
		return items[0]
	}

	return andTerm(items)
}

func UnescapeValue(s string) (string, error) {
	if !strings.ContainsAny(s, `\,=`) {
		return s, nil
	}

	v := bytes.NewBuffer(make([]byte, 0, len(s)))
	isSlash := false
	for _, c := range s {
		if isSlash {
			switch c {
			case '\\', ',', '=':
				v.WriteRune(c)
			default:
				return "", InvalidEscapeSequence{sequence: string([]rune{'\\', c})}
			}
			isSlash = false

			continue
		}

		switch c {
		case '\\':
			isSlash = true
		case ',', '=':
			return "", UnescapedRune{r: c}
		default:
			v.WriteRune(c)
		}
	}

	if isSlash {
		return "", InvalidEscapeSequence{sequence: "\\"}
	}

	return v.String(), nil
}

func Nothing() Selector {
	return nothingSelector{}
}

func Everything() Selector {
	return andTerm{}
}

func parseSelector(selector string, fn TransformFunc) (Selector, error) {
	parts := splitTerms(selector)
	sort.StringSlice(parts).Sort()

	var items []Selector
	for _, part := range parts {
		if part == "" {
			continue
		}

		lhs, op, rhs, ok := splitTerm(part)
		if !ok {
			return nil, fmt.Errorf("Invalid selector: '%s'; can't understand '%s'", selector, part)
		}

		unescapedRHS, err := UnescapeValue(rhs)
		if err != nil {
			return nil, err
		}

		switch op {
		case notEqualOperator:
			items = append(items, &notHasTerm{field: lhs, value: unescapedRHS})
		case doubleEqualOperator:
			items = append(items, &hasTerm{field: lhs, value: unescapedRHS})
		case equalOperator:
			items = append(items, &hasTerm{field: lhs, value: unescapedRHS})
		default:
			return nil, fmt.Errorf("Invalid selector: '%s'; can't understand '%s'", selector, part)
		}
	}

	if len(items) == 1 {
		return items[0].Transform(fn)
	}

	return andTerm(items).Transform(fn)
}

func ParseSelector(selector string) (Selector, error) {
	return parseSelector(selector, func(lhs, rhs string) (newField string, newValue string, err error) {
		return lhs, rhs, nil
	})
}

func ParseSelectorOrDie(s string) Selector {
	selector, err := ParseSelector(s)
	if err != nil {
		panic(err)
	}

	return selector
}

func ParseAndTransformSelector(selector string, fn TransformFunc) (Selector, error) {
	return parseSelector(selector, fn)
}

func OneTermEqualSelector(k, v string) Selector {
	return &hasTerm{field: k, value: v}
}

func OneTermNotEqualSelector(k, v string) Selector {
	return &notHasTerm{field: k, value: v}
}

func AndSelectors(selectors ...Selector) Selector {
	return andTerm(selectors)
}
