package options

import (
	"encoding/json"

	genericOptions "github.com/pachirode/iam_study/internal/pkg/options"
	"github.com/pachirode/iam_study/internal/pkg/server"
	"github.com/pachirode/iam_study/pkg/flags"
	"github.com/pachirode/iam_study/pkg/log"
)

type Options struct {
	GenericServerRunOptions *genericOptions.ServerRunOptions `json:"server" mapstructure:"server"`
	MySQLOptions            *genericOptions.MySQLOptions     `json:"mysql" mapstructure:"mysql"`
	Log                     *log.Options                     `json:"log" mapstructure:"log"`
}

func NewOptions() *Options {
	return &Options{
		GenericServerRunOptions: genericOptions.NewServerRunOptions(),
		MySQLOptions:            genericOptions.NewMySQLOptions(),
		Log:                     log.NewOptions(),
	}
}

func (opt *Options) ApplyTo(config *server.Config) error {
	return nil
}

func (opt *Options) Flags() (nfs flags.NamedFlagSets) {
	opt.GenericServerRunOptions.AddFlags(nfs.GetFlagSet("generic"))
	opt.MySQLOptions.AddFlags(nfs.GetFlagSet("mysql"))
	opt.Log.AddFlags(nfs.GetFlagSet("log"))

	return nfs
}

func (opt *Options) String() string {
	data, _ := json.Marshal(opt)

	return string(data)
}

func (opt *Options) Complete() error {
	return nil
}
