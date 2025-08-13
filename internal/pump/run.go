package pump

import (
	genericAPIServer "github.com/pachirode/iam_study/internal/pkg/server"
	"github.com/pachirode/iam_study/internal/pump/config"
)

func Run(cfg *config.Config, stopCh <-chan struct{}) error {
	go genericAPIServer.ServeHealthCheck(cfg.HealthCheckPath, cfg.HealthCheckAddress)

	server, err := createPumpServer(cfg)
	if err != nil {
		return err
	}

	return server.PrepareRun().Run(stopCh)
}
