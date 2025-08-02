package v1

import "github.com/pachirode/iam_study/internal/apiserver/store"

//go:generate mockgen -self_package=github.com/pachirode/iam_study/internal/apiserver/service/v1 -destination mocke_service.go -package v1 github.com/pachirode/iam_study/internal/apiserver/service/v1 Service,UserSrv

type Service interface {
	Users() UserSrv
}

type service struct {
	store store.Factory
}
