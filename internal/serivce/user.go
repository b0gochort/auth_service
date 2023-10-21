package serivce

import (
	"crypto/sha512"
	"fmt"
	"log/slog"
	"time"
	"unicode"

	apidb "github.com/b0gochort/internal/api_db"
	"github.com/b0gochort/internal/model"
	"github.com/b0gochort/pkg"
)

type UserServiceImpl struct {
	userApiDb apidb.UserApi
}

func NewUserService(userApiDB apidb.UserApi) *UserServiceImpl {
	return &UserServiceImpl{
		userApiDb: userApiDB,
	}
}

func (s *UserServiceImpl) SignUp(userReq model.User) (model.Auth, error) {
	if !verifyPassword(userReq.Password) {
		slog.Info("userService.SignUp.CreateUser: invalid password")
		return model.Auth{}, fmt.Errorf("userService.SignUp :invalid password")
	}

	user := model.UserItem{
		Login:    userReq.Login,
		Password: generatePasswordHash(userReq.Password),
	}

	userId, err := s.userApiDb.CreateUser(user)
	if err != nil {
		slog.Info("userService.SignUp.CreateUser: %s", err.Error())
		return model.Auth{}, err
	}
	token, err := pkg.GenerateToken(userReq.Login, userId, time.Hour*1)
	if err != nil {
		slog.Info("userService.SignUp.GenerateToken: %s", err.Error())
		return model.Auth{}, fmt.Errorf("userService.SignUp.GenerateToken: %s", err.Error())
	}

	auth := model.Auth{
		AccessToken: token,
		Id:          userId,
		Login:       userReq.Login,
	}

	return auth, nil
}

func (s *UserServiceImpl) Login(model.User) (model.Auth, error) {
	return model.Auth{}, nil
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
