package fields

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
