package serivce

import (
	apidb "github.com/b0gochort/internal/api_db"
	"github.com/b0gochort/internal/model"
)

type UserService interface {
	Login(model.User) (uint64, error)
	SignUp(model.User) (model.Tokens, error)
}

type Service struct {
	UserService
}

func NewService(ApiDB *apidb.ApiDB) *Service {
	return &Service{
		UserService: NewUserService(ApiDB.UserApi),
	}
}
