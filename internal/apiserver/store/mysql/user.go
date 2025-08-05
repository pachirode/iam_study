package mysql

import (
	"context"

	"gorm.io/gorm"

	v1 "github.com/pachirode/iam_study/internal/pkg/api/apiserver/v1"
	"github.com/pachirode/iam_study/internal/pkg/code"
	"github.com/pachirode/iam_study/pkg/errors"
	"github.com/pachirode/iam_study/pkg/fields"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
	"github.com/pachirode/iam_study/pkg/utils/gormutil"
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

func (u *users) List(ctx context.Context, opts metaV1.ListOptions) (*v1.UserList, error) {
	ret := &v1.UserList{}
	ol := gormutil.Unpointer(opts.Offset, opts.Limit)

	selector, _ := fields.ParseSelector(opts.FieldSelector)
	username, _ := selector.RequiresExactMatch("name")
	d := u.db.Where("name like ? and status = 1", "%"+username+"%").
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount)

	return ret, d.Error
}

func (u *users) ListOptional(ctx context.Context, opts metaV1.ListOptions) (*v1.UserList, error) {
	ret := &v1.UserList{}
	ol := gormutil.Unpointer(opts.Offset, opts.Limit)

	where := v1.User{}
	whereNot := v1.User{
		IsAdmin: 0,
	}

	selector, _ := fields.ParseSelector(opts.FieldSelector)
	username, found := selector.RequiresExactMatch("name")
	if found {
		where.Name = username
	}

	d := u.db.Where(where).
		Not(whereNot).
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount)

	return ret, d.Error
}
