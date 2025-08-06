package store

import (
	"context"

	v1 "github.com/pachirode/iam_study/internal/pkg/api/apiserver/v1"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
)

type PolicyStore interface {
	Create(ctx context.Context, policy *v1.Policy, opts metaV1.CreateOptions) error
	Update(ctx context.Context, policy *v1.Policy, opts metaV1.UpdateOptions) error
	Delete(ctx context.Context, username string, name string, opts metaV1.DeleteOptions) error
	DeleteCollection(ctx context.Context, username string, names []string, opts metaV1.DeleteOptions) error
	Get(ctx context.Context, username string, name string, opts metaV1.GetOptions) (*v1.Policy, error)
	List(ctx context.Context, username string, opts metaV1.ListOptions) (*v1.PolicyList, error)
}
