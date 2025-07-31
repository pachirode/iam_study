package apiserver

import (
	"github.com/pachirode/iam_study/internal/apiserver/config"
	"github.com/pachirode/iam_study/internal/apiserver/options"
	"github.com/pachirode/iam_study/pkg/app"
	"github.com/pachirode/iam_study/pkg/log"
)

const commandDesc = `The IAM API server validates and configurets data for the api object`

func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("IAM API Server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.Log)
		defer log.Flush()

		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return nil
		}

		return Run(cfg)
	}
}
