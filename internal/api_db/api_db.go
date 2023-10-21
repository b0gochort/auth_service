package apidb

import (
	"github.com/b0gochort/internal/model"
	"github.com/restream/reindexer/v3"
)

type UserApi interface {
	CreateUser(user model.UserItem) (int64, error)
	GetUser(email, password string) (int64, string, error)
	GetUserByIdAndLogin(id int64, login string) error
	UpdatePasswordByLogin(login string) error
}

type ApiDB struct {
	UserApi
}

func NewAPIDB(db *reindexer.Reindexer) *ApiDB {
	return &ApiDB{
		UserApi: NewUserApi(db),
	}
}
