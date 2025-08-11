package authorize

import (
	"github.com/gin-gonic/gin"
	"github.com/ory/ladon"

	"github.com/pachirode/iam_study/internal/authzserver/authorization"
	"github.com/pachirode/iam_study/internal/pkg/code"
	"github.com/pachirode/iam_study/pkg/core"
	"github.com/pachirode/iam_study/pkg/errors"
)

type AuthzController struct {
	store authorization.PolicyGetter
}

func NewAuthzController(store authorization.PolicyGetter) *AuthzController {
	return &AuthzController{
		store: store,
	}
}

func (a *AuthzController) Authorize(ctx *gin.Context) {
	var r ladon.Request
	if err := ctx.ShouldBind(&r); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	auth := authorization.NewAuthorizer(authorization.NewAuthorization(a.store))
	if r.Context == nil {
		r.Context = ladon.Context{}
	}

	r.Context["username"] = ctx.GetString("username")
	rsp := auth.Authorize(ctx, &r)

	core.WriteResponse(ctx, nil, rsp)
}
