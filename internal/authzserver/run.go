package authzserver

import "github.com/pachirode/iam_study/internal/authzserver/authorization/config"

func Run(cfg *config.Config) error {
	server, err := createAuthzServer(cfg)
	if err != nil {
		return err
	}

	return server.PrepareRun().Run()
}
