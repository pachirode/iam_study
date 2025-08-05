package options

import (
	"github.com/spf13/pflag"

	"github.com/pachirode/iam_study/internal/pkg/server"
)

type ServerRunOptions struct {
	Mode        string   `json:"mode"        mapstructure:"mode"`
	Healthz     bool     `json:"healthz"     mapstructure:"healthz"`
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
}

func NewServerRunOptions() *ServerRunOptions {
	defaults := server.NewConfig()

	return &ServerRunOptions{
		Mode:        defaults.Mode,
		Healthz:     defaults.Healthz,
		Middlewares: defaults.Middlewares,
	}
}

func (opt *ServerRunOptions) ApplyTo(config *server.Config) error {
	opt.Mode = config.Mode
	opt.Healthz = config.Healthz
	opt.Middlewares = config.Middlewares

	return nil
}

func (opt *ServerRunOptions) Validate() []error {
	errors := []error{}

	return errors
}

func (opt *ServerRunOptions) AddFlags(flagSet *pflag.FlagSet) {
	flagSet.StringVar(
		&opt.Mode,
		"server.mode",
		opt.Mode,
		"Start server in specified mode. Support: debug, test, release.",
	)
	flagSet.BoolVar(
		&opt.Healthz,
		"server.healthz",
		opt.Healthz,
		"Add self readiness check and install /healthz router.",
	)
	flagSet.StringSliceVar(
		&opt.Middlewares,
		"server.middlewares",
		opt.Middlewares,
		"List of allowed middlewares for server, if empty use default.",
	)
}
