package user

import (
	"github.com/gin-gonic/gin"

	"github.com/pachirode/iam_study/pkg/core"
	"github.com/pachirode/iam_study/pkg/log"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
)

func (u *UserController) Delete(ctx *gin.Context) {
	log.L(ctx).Info("delete user function called.")

	if err := u.serv.Users().Delete(ctx, ctx.Param("name"), metaV1.DeleteOptions{Unscoped: true}); err != nil {
		core.WriteResponse(ctx, err, nil)

		return
	}

	core.WriteResponse(ctx, nil, nil)
}
