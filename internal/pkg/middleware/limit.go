package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var ErrLimitExceeded = errors.New("Limit exceeded")

func Limit(maxEventsPerSec float64, maxBurstSize int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(maxEventsPerSec), maxBurstSize)

	return func(ctx *gin.Context) {
		if limiter.Allow() {
			ctx.Next()

			return
		}

		_ = ctx.Error(ErrLimitExceeded)
		ctx.AbortWithStatus(429)
	}
}
