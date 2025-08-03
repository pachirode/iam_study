package user

import (
	"github.com/gin-gonic/gin"
	"github.com/pachirode/iam_study/internal/pkg/code"
	"github.com/pachirode/iam_study/pkg/auth"
	"github.com/pachirode/iam_study/pkg/core"
	"github.com/pachirode/iam_study/pkg/errors"
	"github.com/pachirode/iam_study/pkg/log"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"omitempty"`
	NewPassword string `json:"newPassword" binding:"password"`
}

func (u *UserController) ChangePassword(ctx *gin.Context) {
	log.L(ctx).Info("Change password function called")

	var r ChangePasswordRequest

	if err := ctx.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	user, err := u.serv.Users().Get(ctx, ctx.Param("name"), metaV1.GetOptions{})
	if err != nil {
		core.WriteResponse(ctx, err, nil)

		return
	}

	if err := user.Compare(r.OldPassword); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrPasswordIncorrect, err.Error()), nil)

		return
	}

	user.Password, _ = auth.Encrypt(r.NewPassword)
	if err := u.serv.Users().ChangePassword(ctx, user); err != nil {
		core.WriteResponse(ctx, err, nil)

		return
	}

	core.WriteResponse(ctx, nil, nil)
}
