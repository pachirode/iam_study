package policy

import (
	serviceV1 "github.com/pachirode/iam_study/internal/apiserver/service/v1"
	"github.com/pachirode/iam_study/internal/apiserver/store"
)

type PolicyController struct {
	serv serviceV1.Service
}

func NewPolicyController(store store.Factory) *PolicyController {
	return &PolicyController{
		serv: serviceV1.NewService(store),
	}
}
