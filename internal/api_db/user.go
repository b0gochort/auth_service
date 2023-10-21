package apidb

import (
	"github.com/b0gochort/internal/model"
	"github.com/restream/reindexer/v3"
)

type UserApiImpl struct {
	db *reindexer.Reindexer
}

func NewUserApi(db *reindexer.Reindexer) *UserApiImpl {
	return &UserApiImpl{
		db: db,
	}
}

func (a *UserApiImpl) CreateUser(model.UserItem) (uint64, error) {
	reind
	return 0, nil
}

func (a *UserApiImpl) FindUser(email, password string) (model.UserItem, error) {
	return model.UserItem{}, nil
}

func (a *UserApiImpl) FindUserById(id uint) (model.UserItem, error) {
	return model.UserItem{}, nil
}

func (a *UserApiImpl) UpdatePasswordByLogin(login string) error {
	return nil
}
