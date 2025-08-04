package apiserver

import (
	"github.com/pachirode/iam_study/internal/apiserver/config"
	"github.com/pachirode/iam_study/internal/apiserver/store/mysql"
	genericOptions "github.com/pachirode/iam_study/internal/pkg/options"
	genericApiServer "github.com/pachirode/iam_study/internal/pkg/server"
	"github.com/pachirode/iam_study/pkg/log"
	"github.com/pachirode/iam_study/pkg/shutdown"
	"github.com/pachirode/iam_study/pkg/shutdown/shutdown_managers/posixsignal"
)

type apiServer struct {
	gracefulShutdown *shutdown.GracefulShutdown
	genericAPIServer *genericApiServer.GenericAPIServer
}

type preparedAPIServer struct {
	*apiServer
}

type ExtraConfig struct {
	Addr         string
	MaxMsgSize   int
	mysqlOptions *genericOptions.MySQLOptions
}

type completedExtraConfig struct {
	*ExtraConfig
}

func (s *apiServer) PrepareRun() preparedAPIServer {
	initRouter(s.genericAPIServer.Engine)

	s.gracefulShutdown.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		mysqlStore, _ := mysql.GetMySQLFactoryOr(nil)
		if mysqlStore != nil {
			_ = mysqlStore.Close()
		}

		s.genericAPIServer.Close()

		return nil
	}))

	return preparedAPIServer{s}
}

func (s preparedAPIServer) Run() error {
	if err := s.gracefulShutdown.Start(); err != nil {
		log.Fatalf("Start shutdown manager failed: %s", err.Error())
	}

	return s.genericAPIServer.Run()
}

func (config *ExtraConfig) complete() *completedExtraConfig {
	if config.Addr == "" {
		config.Addr = "127.0.0.1:8081"
	}

	return &completedExtraConfig{config}
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	gracefulShutdown := shutdown.New()
	gracefulShutdown.AddShutdownManager(posixsignal.NewPosixSignalManager())

	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	server := &apiServer{
		gracefulShutdown: gracefulShutdown,
		genericAPIServer: genericServer,
	}

	return server, nil
}
