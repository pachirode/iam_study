package user

import (
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/pachirode/iam_study/internal/pkg/api/apiserver/v1"
	"github.com/pachirode/iam_study/internal/pkg/code"
	"github.com/pachirode/iam_study/pkg/auth"
	"github.com/pachirode/iam_study/pkg/core"
	"github.com/pachirode/iam_study/pkg/errors"
	"github.com/pachirode/iam_study/pkg/log"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
)

func (u *UserController) Create(ctx *gin.Context) {
	log.L(ctx).Info("user create function called.")

	var r v1.User

	if err := ctx.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	if errs := r.Validate(); len(errs) != 0 {
		core.WriteResponse(ctx, errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), nil)

		return
	}

	r.Password, _ = auth.Encrypt(r.Password)
	r.Status = 1
	r.LoginAt = time.Now()

	if err := u.serv.Users().Create(ctx, &r, metaV1.CreateOptions{}); err != nil {
		core.WriteResponse(ctx, err, nil)

		return
	}

	core.WriteResponse(ctx, nil, r)
}
