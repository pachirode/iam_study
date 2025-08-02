package fields

import "github.com/pachirode/iam_study/pkg/selection"

type Requirements []Requirement

type Requirement struct {
	Operator selection.Operator
	Field    string
	Value    string
}
