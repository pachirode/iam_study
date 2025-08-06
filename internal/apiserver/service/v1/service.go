package v1

import "github.com/pachirode/iam_study/internal/apiserver/store"

//go:generate mockgen -self_package=github.com/pachirode/iam_study/internal/apiserver/service/v1 -destination mocke_service.go -package v1 github.com/pachirode/iam_study/internal/apiserver/service/v1 Service,UserSrv

type Service interface {
	Users() UserSrv
	Secrets() SecretSrv
	Policies() PolicySrv
}

type service struct {
	store store.Factory
}

func NewService(store store.Factory) Service {
	return &service{
		store: store,
	}
}

func (s *service) Users() UserSrv {
	return newUsers(s)
}

func (s *service) Secrets() SecretSrv {
	return newSecrets(s)
}

func (s *service) Policies() PolicySrv {
	return newPolicies(s)
}
