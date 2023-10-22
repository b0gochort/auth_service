package apidb

import (
	"encoding/json"
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

func (a *UserApiImpl) GetUser(email, password string) (model.UserItem, error) {
	err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{})

	if err != nil && err.Error() != "rq: Namespace is already exists" {
		log.Fatal(err)
	}

	elem, ok := a.db.Query("users").Where("email", reindexer.EQ, email).And().Where("password", reindexer.EQ, password).GetJson()
	if !ok {
		return model.UserItem{}, fmt.Errorf("no users with email: %s", email)
	}

	var user model.UserItem

	if err = json.Unmarshal(elem, &user); err != nil {
		return model.UserItem{}, fmt.Errorf("unmarshall")
	}

	return user, nil
}

func (a *UserApiImpl) UpdateUser(user model.UserItem) error {
	err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{})
	if err != nil {
		log.Fatal(err)
	}

	elem := a.db.Query("users").Where("id", reindexer.EQ, user.ID).And().Update()

	if elem.Error() != nil {
		return fmt.Errorf("faild update", elem.Error())
	}

	return nil

}

func (a *UserApiImpl) GetUserByIdAndLogin(id int64, login string) error {
	err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{})
	if err != nil {
		log.Fatal(err)
	}

	elem, ok := a.db.Query("users").Where("login", reindexer.EQ, login).And().Where("id", reindexer.EQ, id).GetJson()
	if !ok {
		return fmt.Errorf("no users with login: %s", login)
	}

	var user model.UserItem

	if err = json.Unmarshal(elem, &user); err != nil {
		return err
	}

	return nil
}

func (a *UserApiImpl) VerificationCode(email string) (model.EmailItem, error) {
	err := a.db.OpenNamespace("codes", reindexer.DefaultNamespaceOptions(), model.EmailItem{})
	if err != nil {
		log.Fatal(err)
	}

	elem, ok := a.db.Query("codes").Where("email", reindexer.EQ, email).GetJson()
	if !ok {
		return model.EmailItem{}, fmt.Errorf("no users with email: %s", email)
	}

	var verification model.EmailItem

	if err = json.Unmarshal(elem, &verification); err != nil {
		return model.EmailItem{}, err
	}

	return verification, nil
}

func (a *UserApiImpl) CreateVerification(verification model.EmailItem) (int64, error) {
	err := a.db.OpenNamespace("codes", reindexer.DefaultNamespaceOptions(), model.EmailItem{})
	if err != nil {
		log.Fatal(err)
	}

	ok, err := a.db.Insert("codes", &verification, "id=serial()")
	if err != nil {
		return 0, err
	}

	if ok == 0 {
		return 0, fmt.Errorf("nil insert")
	}

	return verification.Id, nil
}

func (a *UserApiImpl) SetAuth(email string) error {
	err := a.db.OpenNamespace("codes", reindexer.DefaultNamespaceOptions(), model.EmailItem{})
	if err != nil {
		log.Fatal(err)
	}

	elems := a.db.Query("codes").Where("email", reindexer.EQ, email).Set("Authenticated", true).Update()

	if elems.Error() != nil {
		return fmt.Errorf("faild update", elems.Error())
	}

	return nil
}

func (a *UserApiImpl) CheckAuth(email string) (model.UserItem, error) {
	err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{})
	if err != nil {
		log.Fatal(err)
	}

	elem, ok := a.db.Query("users").Where("email", reindexer.EQ, email).GetJson()

	if !ok {
		return model.UserItem{}, fmt.Errorf("no users with email: %s", email)
	}

	var user model.UserItem

	if err = json.Unmarshal(elem, &user); err != nil {
		return model.UserItem{}, err
	}

	return user, nil

}

func (a *UserApiImpl) UpdatePasswordByLogin(login string) error {
	return nil
}
