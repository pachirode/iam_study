package apiserver

import (
	"context"

	"github.com/pachirode/iam_study/internal/apiserver/config"
	cacheV1 "github.com/pachirode/iam_study/internal/apiserver/controller/v1/cache"
	"github.com/pachirode/iam_study/internal/apiserver/store"
	"github.com/pachirode/iam_study/internal/apiserver/store/mysql"
	pb "github.com/pachirode/iam_study/internal/pkg/api/proto/apiserver/v1"
	genericOptions "github.com/pachirode/iam_study/internal/pkg/options"
	genericApiServer "github.com/pachirode/iam_study/internal/pkg/server"
	"github.com/pachirode/iam_study/pkg/log"
	"github.com/pachirode/iam_study/pkg/shutdown"
	"github.com/pachirode/iam_study/pkg/shutdown/shutdown_managers/posixsignal"
	"github.com/pachirode/iam_study/pkg/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type apiServer struct {
	gracefulShutdown *shutdown.GracefulShutdown
	redisOptions     *genericOptions.RedisOptions
	gRPCAPIServer    *grpcAPIServer
	genericAPIServer *genericApiServer.GenericAPIServer
}

type preparedAPIServer struct {
	*apiServer
}

type ExtraConfig struct {
	Addr         string
	MaxMsgSize   int
	ServerCert   genericOptions.GeneratableKeyCert
	mysqlOptions *genericOptions.MySQLOptions
}

type completedExtraConfig struct {
	*ExtraConfig
}

func (s *apiServer) PrepareRun() preparedAPIServer {
	initRouter(s.genericAPIServer.Engine)

	s.initRedisStore()

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

func (s *apiServer) initRedisStore() {
	ctx, cancel := context.WithCancel(context.Background())
	s.gracefulShutdown.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		cancel()

		return nil
	}))

	config := &storage.Config{
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

	go storage.ConnectToRedis(ctx, config)
}

func (s preparedAPIServer) Run() error {
	go s.gRPCAPIServer.Run()

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

func (config *completedExtraConfig) New() (*grpcAPIServer, error) {
	creds, err := credentials.NewServerTLSFromFile(config.ServerCert.CertKey.CertFile, config.ServerCert.CertKey.KeyFile)
	if err != nil {
		log.Fatalf("Failed to generate credentials %s", grpc.Creds(creds))
	}

	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(config.MaxMsgSize)}
	grpcServer := grpc.NewServer(opts...)

	storeIns, _ := mysql.GetMySQLFactoryOr(config.mysqlOptions)
	store.SetClient(storeIns)
	cacheIns, err := cacheV1.GetCacheInsOr(storeIns)
	if err != nil {
		log.Fatalf("Failed to get cache instance: %s", err.Error())
	}

	pb.RegisterCacheServer(grpcServer, cacheIns)

	reflection.Register(grpcServer)

	return &grpcAPIServer{grpcServer, config.Addr}, nil
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	gracefulShutdown := shutdown.New()
	gracefulShutdown.AddShutdownManager(posixsignal.NewPosixSignalManager())

	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	extraConfig, err := buildExtraConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	extraServer, err := extraConfig.complete().New()
	if err != nil {
		return nil, err
	}

	server := &apiServer{
		gracefulShutdown: gracefulShutdown,
		redisOptions:     cfg.RedisOptions,
		genericAPIServer: genericServer,
		gRPCAPIServer:    extraServer,
	}

	return server, nil
}
