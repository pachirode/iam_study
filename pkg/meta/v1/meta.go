package v1

import "time"

type ObjectMetaAccessor interface {
	GetObjectMeta() Object
}

type Object interface {
	GetID() uint64
	SetID(id uint64)
	GetName() string
	SetName(name string)
	GetCreateAt() time.Time
	SetCreateAt(createAt time.Time)
	GetUpdateAt() time.Time
	SetUpdateAt(updateAt time.Time)
}

type ListInterface interface {
	GetTotalCount() int64
	SetTotalCount(count int64)
}

type Type interface {
	GetAPIVersion() string
	SetAPIVersion(version string)
	GetKind() string
	SetKind(kind string)
}

var _ ListInterface = &ListMeta{}

func (meta *ListMeta) GetTotalCount() int64 {
	return meta.TotalCount
}

func (meta *ListMeta) SetTotalCount(count int64) {
	meta.TotalCount = count
}
