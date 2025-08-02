package fields

import (
	"sort"
	"strings"
)

type Fields interface {
	Has(field string) (exists bool)
	Get(field string) (value string)
}

type Set map[string]string

func (ls Set) String() string {
	selector := make([]string, 0, len(ls))

	for key, value := range ls {
		selector = append(selector, key+"="+value)
	}
	sort.StringSlice(selector).Sort()
	return strings.Join(selector, ",")
}

func (ls Set) Has(field string) bool {
	_, exists := ls[field]
	return exists
}

func (ls Set) Get(field string) string {
	return ls[field]
}
