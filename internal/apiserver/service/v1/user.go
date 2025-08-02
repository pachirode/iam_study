package v1

import (
	"context"
	"sync"

	"github.com/pachirode/iam_study/internal/apiserver/store"
	v1 "github.com/pachirode/iam_study/internal/pkg/api/apiserver/v1"
	"github.com/pachirode/iam_study/internal/pkg/code"
	"github.com/pachirode/iam_study/pkg/errors"
	"github.com/pachirode/iam_study/pkg/log"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
)

type UserSrv interface {
	Creat(ctx context.Context, user *v1.User, opts metaV1.CreateOptions) error
	Update(ctx context.Context, user *v1.User, opts metaV1.UpdateOptions) error
	Delete(ctx context.Context, username string, opts metaV1.DeleteOptions) error
	DeleteCollection(ctx context.Context, usernames []string, opts metaV1.DeleteOptions) error
	Get(ctx context.Context, username string, opts metaV1.GetOptions) (*v1.User, error)
	List(ctx context.Context, opts metaV1.ListOptions) (*v1.UserList, error)
	ListWithBadPerformance(ctx context.Context, opts metaV1.ListOptions) (*v1.UserList, error)
	ChangePassword(ctx context.Context, user *v1.User) error
}

type userService struct {
	store store.Factory
}

func newUsers(srv *service) *userService {
	return &userService{store: srv.store}
}

func (u *userService) List(ctx context.Context, opts metaV1.ListOptions) (*v1.UserList, error) {
	users, err := u.store.Users().List(ctx, opts)
	if err != nil {
		log.L(ctx).Errorf("List users from storage failed: %s", err.Error())

		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	var m sync.Map

	for _, user := range users.Items {
		wg.Add(1)

		go func(user *v1.User) {
			defer wg.Done()

			m.Store(user.ID, &v1.User{
				ObjectMeta: metaV1.ObjectMeta{
					ID:         user.ID,
					InstanceID: user.InstanceID,
					Name:       user.Name,
					Extend:     user.Extend,
					CreatedAt:  user.CreatedAt,
					UpdatedAt:  user.UpdatedAt,
				},
				Nickname:    user.Nickname,
				Email:       user.Email,
				Phone:       user.Phone,
				TotalPolicy: user.TotalPolicy,
				LoginAt:     user.LoginAt,
			})
		}(user)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return nil, err
	}

	infos := make([]*v1.User, 0, len(users.Items))
	for _, user := range users.Items {
		info, _ := m.Load(user.ID)
		infos = append(infos, info.(*v1.User))
	}

	log.L(ctx).Debugf("Get %d users from backend storage", len(infos))

	return &v1.UserList{ListMeta: users.ListMeta, Items: infos}, nil
}
