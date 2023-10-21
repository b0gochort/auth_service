package serivce

import (
	apidb "github.com/b0gochort/internal/api_db"
	"github.com/b0gochort/internal/model"
)

type UserService interface {
	SignUp(user model.User) (uint64, error)
	LogIn(user model.User) (model.Token, error)
}

type Service struct {
	UserService
}

func NewService(ApiDB *apidb.ApiDB) *Service {
	return &Service{
		UserService: NewUserService(ApiDB.UserApi),
	}
}
