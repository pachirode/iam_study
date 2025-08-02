package v1

import (
	"time"

	"github.com/pachirode/iam_study/pkg/scheme"
)

type ObjectMetaAccessor interface {
	GetObjectMeta() Object
}

type Object interface {
	GetID() uint64
	SetID(id uint64)
	GetName() string
	SetName(name string)
	GetCreatedAt() time.Time
	SetCreatedAt(createdAt time.Time)
	GetUpdatedAt() time.Time
	SetUpdatedAt(updatedAt time.Time)
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

func (meta *ListMeta) GetListMeta() ListInterface {
	return meta
}

var _ Type = &TypeMeta{}

func (meta *TypeMeta) GetObjectKind() scheme.ObjectKind {
	return meta
}

func (meta *TypeMeta) SetGroupVersionKind(gvk scheme.GroupVersionKind) {
	meta.APIVersion, meta.Kind = gvk.ToAPIVersionAndKind()
}

func (meta *TypeMeta) GetGroupVersionKind() scheme.GroupVersionKind {
	return scheme.FormAPIVersionAndKind(meta.APIVersion, meta.Kind)
}

func (meta *TypeMeta) GetAPIVersion() string {
	return meta.APIVersion
}

func (meta *TypeMeta) SetAPIVersion(version string) {
	meta.APIVersion = version
}

func (meta *TypeMeta) GetKind() string {
	return meta.Kind
}

func (meta *TypeMeta) SetKind(kind string) {
	meta.Kind = kind
}

var _ Object = &ObjectMeta{}

func (meta *ObjectMeta) GetID() uint64 {
	return meta.ID
}

func (meta *ObjectMeta) SetID(id uint64) {
	meta.ID = id
}

func (meta *ObjectMeta) GetName() string {
	return meta.Name
}

func (meta *ObjectMeta) SetName(name string) {
	meta.Name = name
}

func (meta *ObjectMeta) GetCreatedAt() time.Time {
	return meta.CreatedAt
}

func (meta *ObjectMeta) SetCreatedAt(createdAt time.Time) {
	meta.CreatedAt = createdAt
}

func (meta *ObjectMeta) GetUpdatedAt() time.Time {
	return meta.UpdatedAt
}

func (meta *ObjectMeta) SetUpdatedAt(updatedAt time.Time) {
	meta.UpdatedAt = updatedAt
}

func (meta *ObjectMeta) GetObjectMeta() Object { return meta }
