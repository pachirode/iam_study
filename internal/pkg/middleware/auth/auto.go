package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pachirode/iam_study/internal/pkg/code"
	"github.com/pachirode/iam_study/internal/pkg/middleware"
	"github.com/pachirode/iam_study/pkg/core"
	"github.com/pachirode/iam_study/pkg/errors"
)

const authHeaderCount = 2

type AutoStrategy struct {
	basic BasicStrategy
	jwt   JWTStrategy
}

var _ middleware.AuthStrategy = &AutoStrategy{}

func NewAutoStrategy(basic BasicStrategy, jwt JWTStrategy) AutoStrategy {
	return AutoStrategy{
		basic: basic,
		jwt:   jwt,
	}
}

func (a AutoStrategy) AuthFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		operator := middleware.AuthOperator{}
		authHeader := strings.SplitN(ctx.Request.Header.Get("Authorization"), " ", 2)

		if len(authHeader) != authHeaderCount {
			core.WriteResponse(ctx, errors.WithCode(code.ErrInvalidAuthHeader, "Authorization header format is wrong."), nil)
			ctx.Abort()

			return
		}

		switch authHeader[0] {
		case "Basic":
			operator.SetStrategy(a.basic)
		case "Bearer":
			operator.SetStrategy(a.jwt)
		default:
			core.WriteResponse(ctx, errors.WithCode(code.ErrInvalidAuthHeader, "Unrecognized Authorization header."), nil)
			ctx.Abort()
			return
		}

		operator.AuthFunc()(ctx)

		ctx.Next()
	}
}
