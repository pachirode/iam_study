package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pachirode/iam_study/internal/apiserver/store"
	"github.com/pachirode/iam_study/internal/pkg/code"
	"github.com/pachirode/iam_study/pkg/core"
	"github.com/pachirode/iam_study/pkg/errors"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
)

func Validation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := isAdmin(ctx); err != nil {
			switch ctx.FullPath() {
			case "/v1/users":
				if ctx.Request.Method != http.MethodPost {
					core.WriteResponse(ctx, errors.WithCode(code.ErrPermissionDenied, ""), nil)
					ctx.Abort()

					return
				}
			case "/v1/users/:name", "/v1/users/:name/change_password":
				username := ctx.GetString("username")
				if ctx.Request.Method == http.MethodDelete || (ctx.Request.Method != http.MethodDelete && username != ctx.Param("name")) {
					core.WriteResponse(ctx, errors.WithCode(code.ErrPermissionDenied, ""), nil)
					ctx.Abort()

					return
				}
			default:
			}
		}

		ctx.Next()
	}
}

func isAdmin(ctx *gin.Context) error {
	username := ctx.GetString(UsernameKey)
	user, err := store.Client().Users().Get(ctx, username, metaV1.GetOptions{})
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	if user.IsAdmin != 1 {
		return errors.WithCode(code.ErrPermissionDenied, "user %s is not a administrator", username)
	}

	return nil
}
