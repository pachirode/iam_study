package v1

import (
	"context"

	"github.com/pachirode/iam_study/internal/apiserver/store"
	v1 "github.com/pachirode/iam_study/internal/pkg/api/apiserver/v1"
	"github.com/pachirode/iam_study/internal/pkg/code"
	"github.com/pachirode/iam_study/pkg/errors"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
)

type PolicySrv interface {
	Create(ctx context.Context, policy *v1.Policy, opts metaV1.CreateOptions) error
	Update(ctx context.Context, policy *v1.Policy, opts metaV1.UpdateOptions) error
	Delete(ctx context.Context, username string, name string, opts metaV1.DeleteOptions) error
	DeleteCollection(ctx context.Context, username string, names []string, opts metaV1.DeleteOptions) error
	Get(ctx context.Context, username string, name string, opts metaV1.GetOptions) (*v1.Policy, error)
	List(ctx context.Context, username string, opts metaV1.ListOptions) (*v1.PolicyList, error)
}

type policyService struct {
	store store.Factory
}

var _ PolicySrv = (*policyService)(nil)

func newPolicies(srv *service) *policyService {
	return &policyService{store: srv.store}
}

func (s *policyService) Create(ctx context.Context, policy *v1.Policy, opts metaV1.CreateOptions) error {
	if err := s.store.Policies().Create(ctx, policy, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *policyService) Update(ctx context.Context, policy *v1.Policy, opts metaV1.UpdateOptions) error {
	if err := s.store.Policies().Update(ctx, policy, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *policyService) Delete(ctx context.Context, username, name string, opts metaV1.DeleteOptions) error {
	if err := s.store.Policies().Delete(ctx, username, name, opts); err != nil {
		return err
	}

	return nil
}

func (s *policyService) DeleteCollection(
	ctx context.Context,
	username string,
	names []string,
	opts metaV1.DeleteOptions,
) error {
	if err := s.store.Policies().DeleteCollection(ctx, username, names, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *policyService) Get(ctx context.Context, username, name string, opts metaV1.GetOptions) (*v1.Policy, error) {
	policy, err := s.store.Policies().Get(ctx, username, name, opts)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (s *policyService) List(ctx context.Context, username string, opts metaV1.ListOptions) (*v1.PolicyList, error) {
	policies, err := s.store.Policies().List(ctx, username, opts)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return policies, nil
}
