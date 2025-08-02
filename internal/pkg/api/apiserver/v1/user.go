package v1

import (
	"fmt"
	"time"

	"github.com/pachirode/iam_study/pkg/auth"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
	"github.com/pachirode/iam_study/pkg/utils/idutil"
	"gorm.io/gorm"
)

type User struct {
	metaV1.ObjectMeta `json:"metaData,omitempty"`
	Status            int       `json:"status" gorm:"column:status" validate:"omitempty"`
	Nickname          string    `json:"nickname" gorm:"column:nickname" validate:"required,min=1,max=30"`
	Password          string    `json:"password,omitempty" gorm:"column:password" validate:"required"`
	Email             string    `json:"email" gorm:"column:email" validate:"required,email,min=1,max=100"`
	Phone             string    `json:"phone" gorm:"column:phone" validate:"omitempty"`
	IsAdmin           int       `json:"isAdmin,omitempty" gorm:"isAdmin" validate:"omitempty"`
	TotalPolicy       int64     `json:"totalPolicy" gorm:"-" validate:"omitempty"`
	LoginAt           time.Time `json:"loginAt,omitempty" gorm:"column:loginAt"`
}

type UserList struct {
	metaV1.ListMeta `json:",inline"`
	Items           []*User `json:"items"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) Compare(pwd string) error {
	if err := auth.Compare(u.Password, pwd); err != nil {
		return fmt.Errorf("failed to compile password: %w", err)
	}

	return nil
}

func (u *User) AfterCreate(gdb *gorm.DB) error {
	u.InstanceID = idutil.GetInstanceID(u.ID, "user-")

	return gdb.Save(u).Error
}
