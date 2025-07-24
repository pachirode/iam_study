package errors

import (
	"reflect"
	"sort"
)

type Empty struct{}

type String map[string]Empty

type sortableSliceOfString []string

func (s sortableSliceOfString) Len() int {
	return len(s)
}

func (s sortableSliceOfString) Less(i, j int) bool {
	return lessString(s[i], s[j])
}

func (s sortableSliceOfString) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func lessString(lhs, rhs string) bool {
	return lhs < rhs
}

func (s String) Len() int {
	return len(s)
}

func NewString(items ...string) String {
	ss := String{}
	ss.Insert(items...)
	return ss
}

func StringKeySet(strMap interface{}) String {
	v := reflect.ValueOf(strMap)
	ret := String{}

	for _, key := range v.MapKeys() {
		ret.Insert(key.Interface().(string))
	}

	return ret
}

func (s String) Insert(items ...string) String {
	for _, item := range items {
		s[item] = Empty{}
	}
	return s
}

func (s String) Delete(items ...string) String {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

func (s String) Has(item string) bool {
	_, contained := s[item]
	return contained
}

func (s String) HasAll(items ...string) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}

	return true
}

func (s String) HasAny(items ...string) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}
	return false
}

func (s String) Difference(other String) String {
	res := NewString()
	for k := range s {
		if !other.Has(k) {
			res.Insert(k)
		}
	}
	return res
}

func (s String) Union(other String) String {
	res := NewString()
	for k := range s {
		res.Insert(k)
	}
	for k := range other {
		res.Insert(k)
	}
	return res
}

func (s String) Intersection(other String) String {
	res := NewString()

	if len(s) < len(other) {
		for k := range s {
			if other.Has(k) {
				res.Insert(k)
			}
		}
	} else {
		for k := range other {
			if s.Has(k) {
				res.Insert(k)
			}
		}
	}

	return res
}

func (s String) IsSuperset(other String) bool {
	for item := range other {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

func (s String) Equal(other String) bool {
	return len(s) == len(other) && s.IsSuperset(other)
}

func (s String) List() []string {
	res := make(sortableSliceOfString, 0, len(s))
	for k := range s {
		res = append(res, k)
	}
	sort.Sort(res)
	return res
}

func (s String) UnsortedList() []string {
	res := make([]string, 0, len(s))
	for k := range s {
		res = append(res, k)
	}
	return res
}

func (s String) PopAny() (string, bool) {
	for k := range s {
		s.Delete(k)
		return k, true
	}
	return "", false
}
