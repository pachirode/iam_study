package user

import (
	"github.com/gin-gonic/gin"
	"github.com/pachirode/iam_study/internal/pkg/code"
	"github.com/pachirode/iam_study/pkg/core"
	"github.com/pachirode/iam_study/pkg/errors"
	"github.com/pachirode/iam_study/pkg/log"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
)

func (u *UserController) List(ctx *gin.Context) {
	log.L(ctx).Info("list user function called.")

	var r metaV1.ListOptions
	if err := ctx.ShouldBindQuery(&r); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	users, err := u.serv.Users().List(ctx, r)
	if err != nil {
		core.WriteResponse(ctx, err, nil)

		return
	}

	core.WriteResponse(ctx, nil, users)
}
