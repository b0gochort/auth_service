package apidb

import (
	"fmt"
	"log"

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

func (a *UserApiImpl) CreateUser(user model.UserItem) (int64, error) {
	err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{})
	if err != nil {
		log.Fatal(err)
	}

	ok, err := a.db.Insert("users", &user, "id=serial()")
	if err != nil {
		return 0, err
	}

	if ok == 0 {
		return 0, fmt.Errorf("nil insert")
	}

	return user.ID, nil
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
