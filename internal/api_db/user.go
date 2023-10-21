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

func (a *UserApiImpl) CreateUser(model.UserItem) (int64, error) {
	// reind
	return 0, nil
}

func (a *UserApiImpl) GetUser(email, password string) (int64, error) {
	return 0, nil
}

func (a *UserApiImpl) GetUserByIdAndLogin(id int64, login string) error {
	return nil
}

func (a *UserApiImpl) UpdatePasswordByLogin(login string) error {
	return nil
}
