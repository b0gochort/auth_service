package serivce

import (
	apidb "github.com/b0gochort/internal/api_db"
	"github.com/b0gochort/internal/model"
)

type UserService interface {
	SignUp(user model.User) (model.Auth, error)
	Login(user model.User) (model.Auth, error)
}

type Service struct {
	UserService
}

func NewService(ApiDB *apidb.ApiDB) *Service {
	return &Service{
		UserService: NewUserService(ApiDB.UserApi),
	}
}
