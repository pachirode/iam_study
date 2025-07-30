package server

import (
<<<<<<< HEAD
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pachirode/iam_study/internal/pkg/middleware"
	"github.com/pachirode/iam_study/pkg/log"
=======
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/hcl/hcl/strconv"
	"github.com/pachirode/iam_study/internal/pkg/middleware"
	"github.com/pachirode/iam_study/pkg/core"
	"github.com/pachirode/iam_study/pkg/log"
	"github.com/pachirode/iam_study/pkg/version"
	ginPrometheus "github.com/zsais/go-gin-prometheus"
	"golang.org/x/sync/errgroup"
>>>>>>> 4f132fc (add apiserver run options)
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

<<<<<<< HEAD
	for _, middleware := range server.middlewares {
	}
}

func initGenericAPIServer(server *GenericAPIServer) {
=======
	for _, m := range server.middlewares {
		mv, ok := middleware.Middleware[m]
		if !ok {
			log.Warnf("can not find middleware: %s", m)
			continue
		}

		log.Infof("install middleware: %s", m)
		server.Use(mv)
	}

	server.Use(middleware.Context())
}

func (server *GenericAPIServer) InstallAPIs() {
	if server.healthz {
		server.GET("/healthz", func(ctx *gin.Context) {
			core.WriteResponse(ctx, nil, map[string]string{"status": "ok"})
		})
	}

	if server.enableMetrics {
		prometheus := ginPrometheus.NewPrometheus("gin")
		prometheus.Use(server.Engine)
	}

	if server.enableProfiling {
		pprof.Register(server.Engine)
	}

	server.GET("/version", func(ctx *gin.Context) {
		core.WriteResponse(ctx, nil, version.Get())
	})
}

func initGenericAPIServer(server *GenericAPIServer) {
	server.Setup()
	server.InstallMiddlewares()
	server.InstallAPIs()
}

func (server *GenericAPIServer) Run() error {
	server.insecureServer = &http.Server{
		Addr:    server.InsecureServingInfo.Address,
		Handler: server,
	}

	server.secureServer = &http.Server{
		Addr:    server.SecureServingInfo.Address(),
		Handler: server,
	}

	var eg errgroup.Group

	eg.Go(func() error {
		log.Infof("Start to listening the incoming requests on the http address: %s", server.InsecureServingInfo.Address)

		if err := server.insecureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())

			return err
		}

		log.Infof("Server on %s stopped", server.InsecureServingInfo.Address)
		return nil
	})

	eg.Go(func() error {
		key, cert := server.SecureServingInfo.CertKey.KeyFile, server.SecureServingInfo.CertKey.CertFile
		if cert == "" || key == "" || server.SecureServingInfo.BindPort == 0 {
			return nil
		}

		log.Infof("Start to listening the incoming requests on https address: %s", server.SecureServingInfo.Address())

		if err := server.secureServer.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())

			return err
		}

		log.Infof("Server on %s stopped", server.SecureServingInfo.Address())

		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if server.healthz {
		if err := server.ping(ctx); err != nil {
			return err
		}
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

func (server *GenericAPIServer) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.secureServer.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown secure server failed: %s", err.Error())
	}

	if err := server.insecureServer.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown insecure server failed: %s", err.Error())
	}
}

func (server *GenericAPIServer) ping(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/healthz", server.InsecureServingInfo.Address)

	if strings.Contains(server.InsecureServingInfo.Address, "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%s/healthz", strings.Split(server.InsecureServingInfo.Address, ":")[1])
	}

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil
		}

		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Info("The router has been deployed successfully.")

			resp.Body.Close()

			return nil
		}

		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)

		select {
		case <-ctx.Done():
			log.Fatal("can not ping http server within specified time internal.")
		default:
		}
	}
>>>>>>> 4f132fc (add apiserver run options)
}
