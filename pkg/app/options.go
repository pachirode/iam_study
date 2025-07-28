package app

import "github.com/pachirode/iam_study/pkg/flags"

type ClipOptions interface {
	Flags() (nfs flags.NamedFlagSets)
	Validate() []error
}

type ConfigurableOptions interface {
	ApplyFlags() []error
}

type CompleteableOptions interface {
	Complete() error
}

type PrintableOptions interface {
	String() string
}
