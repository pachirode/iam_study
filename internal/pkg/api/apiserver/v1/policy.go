package v1

import (
	"encoding/json"
	"fmt"

	"github.com/ory/ladon"
	"gorm.io/gorm"

	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
	"github.com/pachirode/iam_study/pkg/utils/idutil"
)

type AuthzPolicy struct {
	ladon.DefaultPolicy
}

type Policy struct {
	metaV1.ObjectMeta `json:"metaData,omitempty"`
	Username          string      `json:"username" gorm:"column:username" validate:"omitempty"`
	Policy            AuthzPolicy `json:"policy,omitempty" gorm:"-" validate:"omitempty"`
	PolicyShadow      string      `json:"-" gorm:"column:policyShadow" validate:"omitempty"`
}

type PolicyList struct {
	metaV1.ListMeta `json:",inline"`
	Items           []*Policy `json:"items"`
}

func (ap AuthzPolicy) String() string {
	data, _ := json.Marshal(ap)

	return string(data)
}

func (p *Policy) TableName() string {
	return "policy"
}

func (p *Policy) BeforeCreate(gdb *gorm.DB) error {
	if err := p.ObjectMeta.BeforeCreate(gdb); err != nil {
		return fmt.Errorf("Failed to run `BeforeCreate` hook: %w", err)
	}

	p.Policy.ID = p.Name
	p.PolicyShadow = p.Policy.String()

	return nil
}

func (p *Policy) AfterCreate(gdb *gorm.DB) error {
	p.InstanceID = idutil.GetInstanceID(p.ID, "policy-")

	return gdb.Save(p).Error
}

func (p *Policy) BeforeUpdate(gdb *gorm.DB) error {
	if err := p.ObjectMeta.BeforeUpdate(gdb); err != nil {
		return fmt.Errorf("Failed to run `BeforeUpdate` hook: %w", err)
	}

	p.Policy.ID = p.Name
	p.PolicyShadow = p.Policy.String()

	return nil
}

func (p *Policy) AfterFind(gdb *gorm.DB) error {
	if err := p.ObjectMeta.AfterFind(gdb); err != nil {
		return fmt.Errorf("Failed to run `AfterFind` hook: %w", err)
	}

	if err := json.Unmarshal([]byte(p.PolicyShadow), &p.Policy); err != nil {
		return fmt.Errorf("Failed to unmarshal policyShadow: %w", err)
	}

	return nil
}
