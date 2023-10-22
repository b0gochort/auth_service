package apidb

import (
	"github.com/b0gochort/internal/model"
	"github.com/restream/reindexer/v3"
)

type UserApi interface {
	CreateUser(user model.UserItem) (int64, error)
	GetUser(email, password string) (model.UserItem, error)
	UpdateUser(user model.UserItem) error
	GetUserByIdAndLogin(id int64, login string) error
	UpdatePasswordByLogin(login string) error
	CreateVerification(verification model.EmailItem) (int64, error)
	VerificationCode(email string) (model.EmailItem, error)
	SetAuth(email string) error
	CheckAuth(email string) (model.UserItem, error)
}

type ApiDB struct {
	UserApi
}

func NewAPIDB(db *reindexer.Reindexer) *ApiDB {
	return &ApiDB{
		UserApi: NewUserApi(db),
	}
}
