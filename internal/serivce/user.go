package serivce

import (
	"crypto/sha512"
	"fmt"
	"unicode"

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

func (s *UserServiceImpl) SignUp(userReq model.User) (uint64, error) {
	if !verifyPassword(userReq.Password) {
		return 0, fmt.Errorf("invalid password")
	}

	user := model.UserItem{
		Login:    userReq.Login,
		Password: generatePasswordHash(userReq.Password),
	}

	userId, err := s.userApiDb.CreateUser(user)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (s *UserServiceImpl) LogIn(model.User) (model.Token, error) {
	return model.Token{}, nil
}

func generatePasswordHash(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))
	//TODO: add to config
	return fmt.Sprintf("%x", hash.Sum([]byte("salt")))
}

func verifyPassword(s string) bool {
	var sevenOrMore, number, upper, special bool
	letters := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
			return false
		}
	}
	sevenOrMore = letters >= 7
	return sevenOrMore && number && upper && special
}
