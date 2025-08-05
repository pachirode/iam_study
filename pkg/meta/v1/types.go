package v1

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Extend define a new type to store filed
type Extend map[string]interface{}

func (ext Extend) String() string {
	data, _ := json.Marshal(ext)
	return string(data)
}

func (ext Extend) Merge(extendShadow string) Extend {
	var extend Extend

	_ = json.Unmarshal([]byte(extendShadow), &extend)
	for k, v := range extend {
		if _, ok := ext[k]; !ok {
			ext[k] = v
		}
	}

	return ext
}

type TypeMeta struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}

type ListMeta struct {
	TotalCount int64 `json:"totalCount,omitempty"`
}

type ObjectMeta struct {
	ID           uint64    `json:"id,omitempty"         gorm:"primary_key;AUTO_INCREMENT;column:id"`
	InstanceID   string    `json:"instanceID,omitempty" gorm:"unqiue;column:instanceID;type:varchar(32);not null"`
	Name         string    `json:"name,omitempty"       gorm:"column:name;type:varchar(64);not null"              validate:"name"`
	Extend       Extend    `json:"extend,omitempty"     gorm:"-"                                                  validate:"omitempty"`
	ExtendShadow string    `json:"-"                    gorm:"column:extendShadow"                                validate:"omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"  gorm:"column:createdAt"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"  gorm:"column:updatedAt"`
}

func (obj *ObjectMeta) BeforeCreate(gdb *gorm.DB) error {
	obj.ExtendShadow = obj.Extend.String()

	return nil
}

func (obj *ObjectMeta) BeforeUpdate(gdb *gorm.DB) error {
	obj.ExtendShadow = obj.Extend.String()

	return nil
}

func (obj *ObjectMeta) AfterFind(gdb *gorm.DB) error {
	if err := json.Unmarshal([]byte(obj.ExtendShadow), &obj.Extend); err != nil {
		return err
	}

	return nil
}

type ListOptions struct {
	TypeMeta       `       json:",inline"`
	LabelSelector  string `json:"labelSelector,omitempty"  form:"labelSelector"`
	FieldSelector  string `json:"fieldSelector,omitempty"  form:"fieldSelector"`
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`
	Offset         *int64 `json:"offset,omitempty"         form:"offset"`
	Limit          *int64 `json:"limit,omitempty"          form:"limit"`
}

type ExportOptions struct {
	TypeMeta `     json:",inline"`
	Export   bool `json:"export"`
	Exact    bool `json:"exact"`
}

type GetOptions struct {
	TypeMeta `json:",inline"`
}

type DeleteOptions struct {
	TypeMeta `     json:",inline"`
	Unscoped bool `json:"unscoped"`
}

type CreateOptions struct {
	TypeMeta `         json:",inline"`
	DryRun   []string `json:"dryRun,omitempty"`
}

type PatchOptions struct {
	TypeMeta `         json:",inline"`
	DryRun   []string `json:"dryRun,omitempty"`
	Force    bool     `json:"force,omitempty"`
}

type UpdateOptions struct {
	TypeMeta `         json:",inline"`
	DryRun   []string `json:"dryRun,omitempty"`
}

type AuthorizeOptions struct {
	TypeMeta `json:",inline"`
}

type TableOptions struct {
	TypeMeta  `     json:",inline"`
	NoHeaders bool `json:"-"`
}
