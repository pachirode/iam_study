package mysql

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type policyAudit struct {
	db *gorm.DB
}

func newPolicyAudits(ds *dataStore) *policyAudit {
	return &policyAudit{ds.db}
}

// ClearOutdated clear data older than a given days.
func (p *policyAudit) ClearOutdated(ctx context.Context, maxReserveDays int) (int64, error) {
	date := time.Now().AddDate(0, 0, -maxReserveDays).Format("2001-01-01 11:00:00")

	d := p.db.Exec("delete from policy_audit where deletedAt < ?", date)

	return d.RowsAffected, d.Error
}
