package apidb

import (
	"github.com/b0gochort/internal/model"
	"github.com/restream/reindexer/v3"
)

type UserApi interface {
	CreateUser(model.UserItem, *reindexer.Reindexer) (uint64, error)
	FindUser(login, password string) (model.UserItem, error)
	FindUserById(id uint) (model.UserItem, error)
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
