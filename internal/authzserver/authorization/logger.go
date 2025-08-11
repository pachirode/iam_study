package authorization

import (
	"context"

	"github.com/ory/ladon"

	"github.com/pachirode/iam_study/pkg/log"
)

type AuditLogger struct {
	client AuthorizationInterface
}

func NewAuditLogger(client AuthorizationInterface) *AuditLogger {
	return &AuditLogger{
		client: client,
	}
}

func (a *AuditLogger) LogRejectedAccessRequest(
	ctx context.Context,
	r *ladon.Request,
	p ladon.Policies,
	d ladon.Policies,
) {
	a.client.LogRejectedAccessRequest(r, p, d)
	log.Debug("Subject access review rejected", log.Any("request", r), log.Any("deciders", d))
}

func (a *AuditLogger) LogGrantedAccessRequest(
	ctx context.Context,
	r *ladon.Request,
	p ladon.Policies,
	d ladon.Policies,
) {
	a.client.LogGrantedAccessRequest(r, p, d)
	log.Debug("Subject access review granted", log.Any("request", r), log.Any("deciders", d))
}
