package policy

import (
	"github.com/gin-gonic/gin"

	v1 "github.com/pachirode/iam_study/internal/pkg/api/apiserver/v1"
	"github.com/pachirode/iam_study/internal/pkg/code"
	"github.com/pachirode/iam_study/internal/pkg/middleware"
	"github.com/pachirode/iam_study/pkg/core"
	"github.com/pachirode/iam_study/pkg/errors"
	"github.com/pachirode/iam_study/pkg/log"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
)

// Update updates policy by the policy identifier.
func (p *PolicyController) Update(c *gin.Context) {
	log.L(c).Info("update policy function called.")

	var r v1.Policy
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	pol, err := p.serv.Policies().Get(c, c.GetString(middleware.UsernameKey), c.Param("name"), metaV1.GetOptions{})
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	// only update policy string
	pol.Policy = r.Policy
	pol.Extend = r.Extend

	if errs := pol.Validate(); len(errs) != 0 {
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), nil)

		return
	}

	if err := p.serv.Policies().Update(c, pol, metaV1.UpdateOptions{}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, pol)
}
