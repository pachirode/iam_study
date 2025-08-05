package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/pachirode/iam_study/pkg/log"
)

const UsernameKey = "username"

func Context() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(log.KeyRequestID, ctx.GetString(XRequestIDKey))
		ctx.Set(log.KeyUsername, ctx.GetString(UsernameKey))
		ctx.Next()
	}
}
