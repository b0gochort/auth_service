package serivce

import (
	"crypto/sha512"
	"fmt"
	"log/slog"
	"net/mail"
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

	if !verifyEmail(userReq.Email) {
		slog.Info("userService.SignUp.CreateUser: invalid email")
		return model.Auth{}, fmt.Errorf("userService.SignUp :invalid email")
	}

	user := model.UserItem{
		Name:       userReq.Name,
		Surname:    userReq.Surname,
		Patronymic: userReq.Patronymic,
		Login:      userReq.Login,
		Password:   generatePasswordHash(userReq.Password),
		Email:      userReq.Email,
		IP:         userReq.IP,
		Birthday:   userReq.Birthday,
		City:       userReq.City,
		Date:       userReq.Date,
		Position:   userReq.Position,
	}

	userId, err := s.userApiDb.CreateUser(user)
	if err != nil {
		slog.Info("userService.SignUp.CreateUser: %s", err.Error())
		return model.Auth{}, fmt.Errorf("userService.SignUp.CreateUser: %s", err.Error())
	}
	token, err := pkg.GenerateToken(userReq.Login, userId, time.Hour*8)
	if err != nil {
		slog.Info("userService.SignUp.GenerateToken: %s", err.Error())
		return model.Auth{}, fmt.Errorf("userService.SignUp.GenerateToken: %s", err.Error())
	}

	return model.Auth{
		AccessToken: token,
		Id:          userId,
		Login:       userReq.Login,
	}, nil
}

func (s *UserServiceImpl) FindUser(userReq model.User) (model.Auth, error) {
	userId, login, err := s.userApiDb.GetUser(userReq.Email, generatePasswordHash(userReq.Password))
	if err != nil {
		slog.Info("userService.Login.FindUser: %s", err.Error())
		return model.Auth{}, err
	}

	token, err := pkg.GenerateToken(userReq.Login, userId, time.Hour*8)
	if err != nil {
		slog.Info("userService.Login.GenerateToken: %s", err.Error())
		return model.Auth{}, fmt.Errorf("userService.Login.GenerateToken: %s", err.Error())
	}

	return model.Auth{
		AccessToken: token,
		Id:          userId,
		Login:       login,
	}, nil
}

func (s *UserServiceImpl) UserExists(userId int64, login string) error {
	if err := s.userApiDb.GetUserByIdAndLogin(userId, login); err != nil {
		slog.Info("userService.AuthMiddleware.FindUserByIdAndLogin: %s", err.Error())
		return fmt.Errorf("userService.Login.AuthMiddleware: %s", err.Error())
	}

	return nil
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
			letters++
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
			letters++

		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
			return false
		}
	}
	sevenOrMore = letters >= 7
	return sevenOrMore && number && upper && special
}

func verifyEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}
