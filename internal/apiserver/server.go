package apiserver

import (
	genericOptions "github.com/pachirode/iam_study/internal/pkg/options"
	genericApiServer "github.com/pachirode/iam_study/internal/pkg/server"
	"github.com/pachirode/iam_study/pkg/shutdown"
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
