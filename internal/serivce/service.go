package serivce

import (
	apidb "github.com/b0gochort/internal/api_db"
	"github.com/b0gochort/internal/model"
)

type UserService interface {
	SignUp(userReq model.User) (model.Auth, error)
	FindUser(userReq model.User) (model.Auth, error)
	UserExists(userId int64, login string) error
}

type Service struct {
	UserService
}

func NewService(ApiDB *apidb.ApiDB) *Service {
	return &Service{
		UserService: NewUserService(ApiDB.UserApi),
	}
}
