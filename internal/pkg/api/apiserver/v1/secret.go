package v1

import (
	"gorm.io/gorm"

	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
	"github.com/pachirode/iam_study/pkg/utils/idutil"
)

type Secret struct {
	metaV1.ObjectMeta `json:"metaData,omitempty"`
	Username          string `json:"username" gorm:"column:username" validate:"omitempty"`
	SecretID          string `json:"secretID" gorm:"column:secretID" validate:"omitempty"`
	SecretKey         string `json:"secretKey" gorm:"column:secretKey" validate:"omitempty"`
	Expires           int64  `json:"expires" gorm:"column:expires" validate:"omitempty"`
	Description       string `json:"description" gorm:"column:description" validate:"description"`
}

type SecretList struct {
	metaV1.ListMeta `json:",inline"`
	Items           []*Secret `json:"items"`
}

func (s *Secret) TableName() string {
	return "secret"
}

func (s *Secret) AfterCreate(gdb *gorm.DB) error {
	s.InstanceID = idutil.GetInstanceID(s.ID, "secret-")

	return gdb.Save(s).Error
}
