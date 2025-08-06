package store

import (
	"context"

	v1 "github.com/pachirode/iam_study/internal/pkg/api/apiserver/v1"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
)

type SecretStore interface {
	Create(ctx context.Context, secret *v1.Secret, opts metaV1.CreateOptions) error
	Update(ctx context.Context, secret *v1.Secret, opts metaV1.UpdateOptions) error
	Delete(ctx context.Context, username, secretID string, opts metaV1.DeleteOptions) error
	DeleteCollection(ctx context.Context, username string, secretIDs []string, opts metaV1.DeleteOptions) error
	Get(ctx context.Context, username, secretID string, opts metaV1.GetOptions) (*v1.Secret, error)
	List(ctx context.Context, username string, opts metaV1.ListOptions) (*v1.SecretList, error)
}
