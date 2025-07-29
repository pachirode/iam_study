package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	XRequestIDKey = "X-Request-ID"
)

func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader(XRequestIDKey)

		if requestID == "" {
			requestID = uuid.NewString()
			ctx.Request.Header.Set(XRequestIDKey, requestID)
			ctx.Set(XRequestIDKey, requestID)
		}

		ctx.Writer.Header().Set(XRequestIDKey, requestID)
		ctx.Next()
	}
}
