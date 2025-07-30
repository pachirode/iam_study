package config

import "github.com/pachirode/iam_study/internal/apiserver/options"

type Config struct {
	*options.Options
}

func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{
		opts,
	}, nil
}
