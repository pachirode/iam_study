package apiserver

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/pachirode/iam_study/internal/apiserver/config"
	v1 "github.com/pachirode/iam_study/internal/pkg/api/apiserver/v1"
	genericApiServer "github.com/pachirode/iam_study/internal/pkg/server"
	"github.com/pachirode/iam_study/pkg/log"
)

func parseWithHeader(ctx *gin.Context) (loginInfo, error) {
	auth := strings.SplitN(ctx.Request.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		log.Errorf("Get basic string from Authorization header failed")

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	payload, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		log.Errorf("Decode basic string: %s", err.Error())

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		log.Errorf("Parse payload failed")

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	return loginInfo{
		Username: pair[0],
		Password: pair[1],
	}, nil
}

func parseWithBody(ctx *gin.Context) (loginInfo, error) {
	var login loginInfo

	if err := ctx.ShouldBindJSON(&login); err != nil {
		log.Errorf("Parse login parameters: %s", err.Error())

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	return login, nil
}

func refreshResponse() func(ctx *gin.Context, code int, token string, expire time.Time) {
	return func(ctx *gin.Context, code int, token string, expire time.Time) {
		ctx.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func loginResponse() func(ctx *gin.Context, code int, token string, expire time.Time) {
	return func(ctx *gin.Context, code int, token string, expire time.Time) {
		ctx.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		claims := jwt.MapClaims{
			"iss": APIServerIssUser,
			"aud": APIServerAudience,
		}

		if u, ok := data.(*v1.User); ok {
			claims[jwt.IdentityKey] = u.Name
			claims["sub"] = u.Name
		}

		return claims
	}
}

func authorizator() func(data interface{}, ctx *gin.Context) bool {
	return func(data interface{}, ctx *gin.Context) bool {
		if v, ok := data.(string); ok {
			log.L(ctx).Infof("user `%s` is authenticated.", v)

			return true
		}

		return false
	}
}

func buildGenericConfig(cfg *config.Config) (genericConfig *genericApiServer.Config, lastErr error) {
	genericConfig = genericApiServer.NewConfig()
	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return
}

func buildExtraConfig(cfg *config.Config) (*ExtraConfig, error) {
	return &ExtraConfig{
		Addr:         fmt.Sprintf("%s:%d", cfg.GRPCOptions.BindAddress, cfg.GRPCOptions.BindPort),
		MaxMsgSize:   cfg.GRPCOptions.MaxMsgSize,
		mysqlOptions: cfg.MySQLOptions,
	}, nil
}
