package scheme

type ObjectKind interface {
	GetGroupVersionKind() GroupVersionKind
	SetGroupVersionKind(groupKind GroupVersionKind)
}

var EmptyObjectKind = emptyObjectKind{}

func (emptyObjectKind) GetGroupVersionKind() GroupVersionKind {
	return GroupVersionKind{}
}

func (emptyObjectKind) SetGroupVersionKind(gvk GroupVersionKind) {
}
