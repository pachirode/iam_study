package authzserver

import (
	"context"

	"github.com/pachirode/iam_study/internal/authzserver/analytics"
	"github.com/pachirode/iam_study/internal/authzserver/authorization/config"
	"github.com/pachirode/iam_study/internal/authzserver/load"
	"github.com/pachirode/iam_study/internal/authzserver/load/cache"
	"github.com/pachirode/iam_study/internal/authzserver/store/apiserver"
	genericOptions "github.com/pachirode/iam_study/internal/pkg/options"
	genericAPIServer "github.com/pachirode/iam_study/internal/pkg/server"
	"github.com/pachirode/iam_study/pkg/errors"
	"github.com/pachirode/iam_study/pkg/log"
	"github.com/pachirode/iam_study/pkg/shutdown"
	"github.com/pachirode/iam_study/pkg/shutdown/shutdown_managers/posixsignal"
	"github.com/pachirode/iam_study/pkg/storage"
)

const RedisKeyPrefix = "analytics-"

type authzServer struct {
	gracefulShudown  *shutdown.GracefulShutdown
	rpcServer        string
	clientCA         string
	redisOptions     *genericOptions.RedisOptions
	genericAPIServer *genericAPIServer.GenericAPIServer
	analyticsOptions *analytics.AnalyticsOptions
	redisCancelFunc  context.CancelFunc
}

type preparedAuthzServer struct {
	*authzServer
}

func createAuthzServer(cfg *config.Config) (*authzServer, error) {
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

	server := &authzServer{
		gracefulShudown:  gracefulShutdown,
		redisOptions:     cfg.RedisOptions,
		analyticsOptions: cfg.AnalyticsOptions,
		rpcServer:        cfg.RPCServer,
		clientCA:         cfg.ClientCA,
		genericAPIServer: genericServer,
	}

	return server, nil
}

func (s *authzServer) PrepareRun() preparedAuthzServer {
	_ = s.initialize()

	initRouter(s.genericAPIServer.Engine)

	return preparedAuthzServer{s}
}

func (s preparedAuthzServer) Run() error {
	s.gracefulShudown.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		s.genericAPIServer.Close()
		if s.analyticsOptions.Enable {
			analytics.GetAnalytics().Stop()
		}
		s.redisCancelFunc()

		return nil
	}))

	if err := s.gracefulShudown.Start(); err != nil {
		log.Fatalf("Start shutdown manager failed: %s", err.Error())
	}

	return s.genericAPIServer.Run()
}

func buildGenericConfig(cfg *config.Config) (genericConfig *genericAPIServer.Config, lastErr error) {
	genericConfig = genericAPIServer.NewConfig()

	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return
}

func (s *authzServer) buildStorageConfig() *storage.Config {
	return &storage.Config{
		Host:                  s.redisOptions.Host,
		Port:                  s.redisOptions.Port,
		Addrs:                 s.redisOptions.Addrs,
		MasterName:            s.redisOptions.MasterName,
		Username:              s.redisOptions.Username,
		Password:              s.redisOptions.Password,
		Database:              s.redisOptions.Database,
		MaxIdle:               s.redisOptions.MaxIdle,
		MaxActive:             s.redisOptions.MaxActive,
		Timeout:               s.redisOptions.Timeout,
		EnableCluster:         s.redisOptions.EnableCluster,
		UseSSL:                s.redisOptions.UseSSL,
		SSLInsecureSkipVerify: s.redisOptions.SSLInsecureSkipVerify,
	}
}

func (s *authzServer) initialize() error {
	ctx, cancel := context.WithCancel(context.Background())
	s.redisCancelFunc = cancel

	go storage.ConnectToRedis(ctx, s.buildStorageConfig())

	cacheIns, err := cache.GetCacheInsOr(apiserver.GetAPIServerFactoryOrDie(s.rpcServer, s.clientCA))
	if err != nil {
		return errors.Wrap(err, "get cache instance failed")
	}

	load.NewLoader(ctx, cacheIns).Start()

	if s.analyticsOptions.Enable {
		analyticsStore := storage.RedisCluster{KeyPrefix: RedisKeyPrefix}
		analyticsIns := analytics.NewAnalytics(s.analyticsOptions, &analyticsStore)
		analyticsIns.Start()
	}

	return nil
}
