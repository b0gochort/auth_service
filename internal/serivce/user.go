package serivce

import (
	apidb "github.com/b0gochort/internal/api_db"
	"github.com/b0gochort/internal/model"
)

type UserServiceImpl struct {
	userApiDb apidb.UserApi
}

func NewUserService(userApiDB apidb.UserApi) *UserServiceImpl {
	return &UserServiceImpl{
		userApiDb: userApiDB,
	}
}

func (s *UserServiceImpl) Login(model.User) (uint64, error) {
	return 0, nil
}

func (s *UserServiceImpl) SignUp(model.User) (model.Tokens, error) {
	return model.Tokens{}, nil
}
