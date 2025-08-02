package sets

type String map[string]Empty

func NewString(items ...string) String {
	strMap := String{}
	strMap.Insert(items...)
	return strMap
}

func (strMap String) Insert(items ...string) String {
	for _, item := range items {
		strMap[item] = Empty{}
	}

	return strMap
}

func (strMap String) Delete(items ...string) String {
	for _, item := range items {
		delete(strMap, item)
	}
	return strMap
}

func (strMap String) Has(item string) bool {
	_, contained := strMap[item]
	return contained
}

func (strMap String) HasAll(items ...string) bool {
	for _, item := range items {
		if !strMap.Has(item) {
			return false
		}
	}
	return true
}

func (strMap String) HasAny(items ...string) bool {
	for _, item := range items {
		if strMap.Has(item) {
			return true
		}
	}
	return false
}

func (strMap String) Difference(other String) String {
	result := NewString()
	for key := range strMap {
		if !other.Has(key) {
			result.Insert(key)
		}
	}

	return result
}

func (strMap String) Union(other String) String {
	result := NewString()
	for key := range strMap {
		result.Insert(key)
	}
	for key := range other {
		result.Insert(key)
	}
	return result
}
