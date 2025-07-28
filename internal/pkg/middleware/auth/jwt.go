package auth

import (
	ginJwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/pachirode/iam_study/internal/pkg/middleware"
)

const (
	AuthzAudience = "iam.authz"
)

var _ middleware.AuthStrategy = &JWTStrategy{}

type JWTStrategy struct {
	ginJwt.GinJWTMiddleware
}

func NewJWTStrategy(gJwt ginJwt.GinJWTMiddleware) JWTStrategy {
	return JWTStrategy{gJwt}
}

func (jwtStrategy JWTStrategy) AuthFunc() gin.HandlerFunc {
	return jwtStrategy.MiddlewareFunc()
}
