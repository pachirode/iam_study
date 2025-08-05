package auth

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"

	"github.com/pachirode/iam_study/internal/pkg/code"
	"github.com/pachirode/iam_study/internal/pkg/middleware"
	"github.com/pachirode/iam_study/pkg/core"
	"github.com/pachirode/iam_study/pkg/errors"
)

var (
	ErrMissingKID    = errors.New("Invalid token format: missing kid field in claims")
	ErrMissingSecret = errors.New("Can not obtain sercret information from cache")
)

type Secret struct {
	Username string
	ID       string
	Key      string
	Expires  int64
}

type CacheStrategy struct {
	get func(kid string) (Secret, error)
}

var _ middleware.AuthStrategy = &CacheStrategy{}

func NewCacheStrategy(get func(kid string) (Secret, error)) CacheStrategy {
	return CacheStrategy{get}
}

func (cache CacheStrategy) AuthFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.Request.Header.Get("Authorization")
		if len(header) == 0 {
			core.WriteResponse(
				ctx,
				errors.WithCode(code.ErrMissingHeader, "Authorization header can not be empty."),
				nil,
			)
			ctx.Abort()

			return
		}

		var (
			rawJWT  string
			sercret Secret
			err     error
		)
		fmt.Sscanf(header, "Bearer %s", &rawJWT)

		claims := &jwt.MapClaims{}

		parsedToken, err := jwt.ParseWithClaims(rawJWT, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			kid, ok := token.Header["kid"].(string)

			if !ok {
				return nil, ErrMissingKID
			}

			sercret, err = cache.get(kid)
			if err != nil {
				return nil, ErrMissingSecret
			}

			return []byte(sercret.Key), nil
		}, jwt.WithAudience(AuthzAudience))

		if err != nil || !parsedToken.Valid {
			core.WriteResponse(ctx, errors.WithCode(code.ErrSignatureInvalid, err.Error()), nil)
			ctx.Abort()

			return
		}

		ctx.Set(middleware.UsernameKey, sercret.Username)
		ctx.Next()
	}
}

func KeyExpired(expires int64) bool {
	if expires >= 1 {
		return time.Now().After(time.Unix(expires, 0))
	}

	return false
}
