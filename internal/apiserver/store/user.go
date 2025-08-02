package store

import (
	"context"

	v1 "github.com/pachirode/iam_study/internal/pkg/api/apiserver/v1"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
)

type UserStore interface {
	Create(ctx context.Context, user *v1.User, opts metaV1.CreateOptions) error
	Update(ctx context.Context, user *v1.User, opts metaV1.UpdateOptions) error
	Delete(ctx context.Context, username string, opts metaV1.DeleteOptions) error
	DeleteCollection(ctx context.Context, usernames []string, opts metaV1.DeleteOptions) error
	Get(ctx context.Context, username string, opts metaV1.GetOptions) (*v1.User, error)
	List(ctx context.Context, opts metaV1.ListOptions) (*v1.UserList, error)
}
