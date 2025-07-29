package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pachirode/iam_study/internal/pkg/middleware"
	"github.com/pachirode/iam_study/pkg/log"
)

type GenericAPIServer struct {
	middlewares         []string
	mode                string
	SecureServingInfo   *SecureServingInfo
	InsecureServingInfo *InsecureServingInfo
	ShutdownTimeout     time.Duration

	*gin.Engine
	healthz         bool
	enableMetrics   bool
	enableProfiling bool

	insecureServer, secureServer *http.Server
}

func (server *GenericAPIServer) Setup() {
	gin.SetMode(server.mode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

func (server *GenericAPIServer) InstallMiddlewares() {
	server.Use(middleware.RequestID())

	for _, middleware := range server.middlewares {
	}
}

func initGenericAPIServer(server *GenericAPIServer) {
}
