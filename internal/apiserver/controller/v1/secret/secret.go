package secret

import (
	serviceV1 "github.com/pachirode/iam_study/internal/apiserver/service/v1"
	"github.com/pachirode/iam_study/internal/apiserver/store"
)

// SecretController create a secret handler used to handle request for secret resource.
type SecretController struct {
	serve serviceV1.Service
}

// NewSecretController creates a secret handler.
func NewSecretController(store store.Factory) *SecretController {
	return &SecretController{
		serve: serviceV1.NewService(store),
	}
}
