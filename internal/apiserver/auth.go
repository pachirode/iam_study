package apiserver

import (
	"context"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/pachirode/iam_study/internal/apiserver/store"
	"github.com/pachirode/iam_study/internal/pkg/middleware"
	"github.com/pachirode/iam_study/internal/pkg/middleware/auth"
	"github.com/pachirode/iam_study/pkg/log"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
	"github.com/spf13/viper"
)

const (
	APIServerAudience = "iam.api.pachirode.com"
	APIServerIssUser  = "iam.apiserver"
)

type loginInfo struct {
	Username string `form:"username" json:"username" binding:"required,username"`
	Password string `form:"password" json:"password" binding:"required,password"`
}

func newBasicAuth() middleware.AuthStrategy {
	return auth.NewBasicStrategy(func(username, password string) bool {
		user, err := store.Client().Users().Get(context.TODO(), username, metaV1.GetOptions{})
		if err != nil {
			return false
		}

		if err := user.Compare(password); err != nil {
			return false
		}

		user.LoginAt = time.Now()
		_ = store.Client().Users().Update(context.TODO(), user, metaV1.UpdateOptions{})

		return true
	})
}

func newJwtAuth() middleware.AuthStrategy {
	ginJwt, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            viper.GetString("jwt.Realm"),
		SigningAlgorithm: "HS256",
		Key:              []byte(viper.GetString("jwt.key")),
		Timeout:          viper.GetDuration("jwt.timeout"),
		MaxRefresh:       viper.GetDuration("jwt.max-refresh"),
		Authenticator:    authenticator(),
		LoginResponse:    loginResponse(),
		LogoutResponse: func(ctx *gin.Context, code int) {
			ctx.JSON(http.StatusOK, nil)
		},
		RefreshResponse: refreshResponse(),
		PayloadFunc:     payloadFunc(),
		IdentityHandler: func(ctx *gin.Context) interface{} {
			claims := jwt.ExtractClaims(ctx)

			return claims[jwt.IdentityKey]
		},
		IdentityKey:  middleware.UsernameKey,
		Authorizator: authorizator(),
		Unauthorized: func(ctx *gin.Context, code int, message string) {
			ctx.JSON(code, gin.H{
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		SendCookie:    true,
		TimeFunc:      time.Now,
	})

	return auth.NewJWTStrategy(*ginJwt)
}

func newAutoAuth() middleware.AuthStrategy {
	return auth.NewAutoStrategy(newBasicAuth().(auth.BasicStrategy), newJwtAuth().(auth.JWTStrategy))
}

func authenticator() func(ctx *gin.Context) (interface{}, error) {
	return func(ctx *gin.Context) (interface{}, error) {
		var login loginInfo
		var err error

		if ctx.Request.Header.Get("Authorization") != "" {
			login, err = parseWithHeader(ctx)
		} else {
			login, err = parseWithBody(ctx)
		}

		if err != nil {
			return "", jwt.ErrFailedAuthentication
		}

		user, err := store.Client().Users().Get(ctx, login.Username, metaV1.GetOptions{})
		if err != nil {
			log.Errorf("Get user information failed: %s", err.Error())

			return "", jwt.ErrFailedAuthentication
		}

		if err := user.Compare(login.Password); err != nil {
			return "", jwt.ErrFailedAuthentication
		}

		user.LoginAt = time.Now()
		_ = store.Client().Users().Update(ctx, user, metaV1.UpdateOptions{})

		return user, nil
	}
}
