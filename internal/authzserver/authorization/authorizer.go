package authorization

import (
	"context"

	"github.com/ory/ladon"

	authzV1 "github.com/pachirode/iam_study/internal/pkg/api/authz/v1"
	"github.com/pachirode/iam_study/pkg/log"
)

type Authorizer struct {
	warden ladon.Warden
}

func NewAuthorizer(authorizationClient AuthorizationInterface) *Authorizer {
	return &Authorizer{
		warden: &ladon.Ladon{
			Manager:     NewPolicyManager(authorizationClient),
			AuditLogger: NewAuditLogger(authorizationClient),
		},
	}
}

func (a *Authorizer) Authorize(ctx context.Context, request *ladon.Request) *authzV1.Response {
	log.Debug("Authorize request", log.Any("request", request))

	if err := a.warden.IsAllowed(ctx, request); err != nil {
		return &authzV1.Response{
			Denied: true,
			Reason: err.Error(),
		}
	}

	return &authzV1.Response{
		Allowed: true,
	}
}
