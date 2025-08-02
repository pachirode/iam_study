package mysql

import (
	"context"
	"os/user"

	v1 "github.com/pachirode/iam_study/internal/pkg/api/apiserver/v1"
	"github.com/pachirode/iam_study/internal/pkg/code"
	"github.com/pachirode/iam_study/pkg/errors"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
	"github.com/pachirode/iam_study/pkg/utils/gormutil"
	"gorm.io/gorm"
)

type users struct {
	db *gorm.DB
}

func newUsers(ds *dataStore) *users {
	return &users{ds.db}
}

func (u *users) Create(ctx context.Context, user *v1.User, opts metaV1.CreateOptions) error {
	return u.db.Create(&user).Error
}

func (u *users) Update(ctx context.Context, user *v1.User, opts metaV1.UpdateOptions) error {
	return u.db.Save(user).Error
}

func (u *users) Delete(ctx context.Context, username string, opts metaV1.DeleteOptions) error {
	err := u.db.Where("name = ?", username).Delete(&v1.User{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (u *users) DeleteCollection(ctx context.Context, usernames []string, opts metaV1.DeleteOptions) error {
	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}

	return u.db.Where("name in (?)", usernames).Delete(&v1.User{}).Error
}

func (u *users) Get(ctx context.Context, username string, opts metaV1.GetOptions) (*v1.User, error) {
	user := &v1.User{}
	err := u.db.Where("name = ? and status = 1", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return user, nil
}
