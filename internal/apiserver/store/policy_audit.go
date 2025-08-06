package store

import (
	"context"
)

type PolicyAuditStore interface {
	ClearOutdated(ctx context.Context, maxReserveDays int) (int64, error)
}
