package apiserver

import (
	"github.com/gin-gonic/gin"

	"github.com/pachirode/iam_study/internal/apiserver/controller/v1/user"
	"github.com/pachirode/iam_study/internal/apiserver/store/mysql"
	"github.com/pachirode/iam_study/internal/pkg/code"
	"github.com/pachirode/iam_study/internal/pkg/middleware"
	"github.com/pachirode/iam_study/internal/pkg/middleware/auth"
	"github.com/pachirode/iam_study/pkg/core"
	"github.com/pachirode/iam_study/pkg/errors"
)

func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine) *gin.Engine {
	jwtStrategy, _ := newJwtAuth().(auth.JWTStrategy)
	g.POST("/login", jwtStrategy.LoginHandler)
	g.POST("/logout", jwtStrategy.LogoutHandler)
	g.POST("/refresh", jwtStrategy.RefreshHandler)

	auto := newAutoAuth()
	g.NoRoute(auto.AuthFunc(), func(ctx *gin.Context) {
		core.WriteResponse(ctx, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	})

	storeIns, _ := mysql.GetMySQLFactoryOr(nil)
	v1 := g.Group("/v1")
	{
		userV1 := v1.Group("/users")
		{
			userController := user.NewUserController(storeIns)

			userV1.POST("", userController.Create)
			userV1.Use(auto.AuthFunc(), middleware.Validation())
			userV1.DELETE("", userController.DeleteCollection)
			userV1.DELETE(":name", userController.Delete)
			userV1.PUT(":name/change-password", userController.ChangePassword)
			userV1.PUT(":name", userController.Update)
			userV1.GET("", userController.List)
			userV1.GET(":name", userController.Get)
		}

		v1.Use(auto.AuthFunc())
	}

	return g
}

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}
