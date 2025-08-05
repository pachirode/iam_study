package user

import (
	"github.com/gin-gonic/gin"

	"github.com/pachirode/iam_study/pkg/core"
	"github.com/pachirode/iam_study/pkg/log"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
)

func (u *UserController) DeleteCollection(ctx *gin.Context) {
	log.L(ctx).Info("batch delete user function called.")

	usernames := ctx.QueryArray("name")

	if err := u.serv.Users().DeleteCollection(ctx, usernames, metaV1.DeleteOptions{}); err != nil {
		core.WriteResponse(ctx, err, nil)

		return
	}

	core.WriteResponse(ctx, nil, nil)
}
