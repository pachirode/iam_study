package authzserver

import (
	"github.com/gin-gonic/gin"

	"github.com/pachirode/iam_study/internal/authzserver/controller/v1/authorize"
	"github.com/pachirode/iam_study/internal/authzserver/load/cache"
	"github.com/pachirode/iam_study/internal/pkg/code"
	"github.com/pachirode/iam_study/pkg/core"
	"github.com/pachirode/iam_study/pkg/errors"
	"github.com/pachirode/iam_study/pkg/log"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine) *gin.Engine {
	auth := newCacheAuth()
	g.NoRoute(auth.AuthFunc(), func(ctx *gin.Context) {
		core.WriteResponse(ctx, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	})

	cacheIns, _ := cache.GetCacheInsOr(nil)
	if cacheIns == nil {
		log.Panicf("get nil cache instance")
	}

	apiV1 := g.Group("/v1", auth.AuthFunc())
	{
		authzController := authorize.NewAuthzController(cacheIns)
		apiV1.POST("/authz", authzController.Authorize)
	}

	return g
}
