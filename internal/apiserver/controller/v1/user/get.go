package user

import (
	"github.com/gin-gonic/gin"
	"github.com/pachirode/iam_study/pkg/core"
	"github.com/pachirode/iam_study/pkg/log"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
)

func (u *UserController) Get(ctx *gin.Context) {
	log.L(ctx).Info("get user function called.")

	user, err := u.serv.Users().Get(ctx, ctx.Param("name"), metaV1.GetOptions{})
	if err != nil {
		core.WriteResponse(ctx, err, nil)

		return
	}

	core.WriteResponse(ctx, nil, user)
}
