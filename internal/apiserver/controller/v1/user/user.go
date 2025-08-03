package user

import (
	serviceV1 "github.com/pachirode/iam_study/internal/apiserver/service/v1"
	"github.com/pachirode/iam_study/internal/apiserver/store"
)

type UserController struct {
	serv serviceV1.Service
}

func NewUserController(store store.Factory) *UserController {
	return &UserController{
		serv: serviceV1.NewService(store),
	}
}
