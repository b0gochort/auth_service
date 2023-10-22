package serivce

import (
	"crypto/sha512"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/mail"
	"time"
	"unicode"

	"github.com/spf13/viper"

	gomail "gopkg.in/mail.v2"

	apidb "github.com/b0gochort/internal/api_db"
	"github.com/b0gochort/internal/model"
	"github.com/b0gochort/pkg"
)

const ConfirmationCodeLifetime = 1 * time.Hour

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
	user, err := s.userApiDb.GetUser(userReq.Email, generatePasswordHash(userReq.Password))
	if err != nil {
		slog.Info("userService.Login.FindUser: %s", err.Error())
		return model.Auth{}, err
	}

	userNew := model.UserItem{
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

	if err := s.userApiDb.UpdateUser(userNew); err != nil {
		slog.Info("userService.Login.UpdateUser: %s", err.Error())
		return model.Auth{}, err
	}

	token, err := pkg.GenerateToken(user.Login, user.ID, time.Hour*8)
	if err != nil {
		slog.Info("userService.Login.GenerateToken: %s", err.Error())
		return model.Auth{}, fmt.Errorf("userService.Login.GenerateToken: %s", err.Error())
	}

	return model.Auth{
		AccessToken: token,
		Id:          user.ID,
		Login:       user.Login,
	}, nil
}

type EmailMessage struct {
	Email   string
	Subject string
	Message string
	Path    string
}

func (s *UserServiceImpl) SendConfirmationEmail(email string) (int64, error) {
	verification := model.EmailItem{
		Email: email,
		Code:  pkg.GenerateCode(),
		Time:  time.Now().Unix(),
	}

	emailMessage := EmailMessage{
		Email:   verification.Email,
		Subject: "Subject: Код подтверждения\n",
		Message: verification.Code,
	}
	err := emailMessage.sendEmail()
	if err != nil {
		slog.Info("userService.SendConfirmationEmail.sendEmal: %s", err.Error())
		return 0, err
	}

	codeId, err := s.userApiDb.CreateVerification(verification)
	if err != nil {
		slog.Info("userService.SendConfirmationEmail.CreateVerification: %s", err.Error())
		return 0, err
	}

	return codeId, nil
}

func (s UserServiceImpl) VerificateEmailCode(code, email string) error {
	res, err := s.userApiDb.VerificationCode(email)
	if err != nil {
		slog.Info("userService.VerificateEmailCode.FindUserByIdAndLogin: %s", err.Error())
		return fmt.Errorf("userService.Login.VerificateEmailCode: %s", err.Error())
	}

	if res.Code != code {
		return fmt.Errorf("invalid code, repeat")
	}

	if err := s.userApiDb.SetAuth(email); err != nil {
		slog.Info("userService.VerificateEmailCode.SetAuth: %s", err.Error())
		return err
	}

	return nil
}

func (s *UserServiceImpl) UserExists(userId int64, login string) error {
	if err := s.userApiDb.GetUserByIdAndLogin(userId, login); err != nil {
		slog.Info("userService.AuthMiddleware.FindUserByIdAndLogin: %s", err.Error())
		return fmt.Errorf("userService.Login.AuthMiddleware: %s", err.Error())
	}

	return nil
}

func (s *UserServiceImpl) CheckUserAuth(email string) error {
	user, err := s.userApiDb.CheckAuth(email)
	if err != nil {
		slog.Info("userService.VerificateEmailCode.ChechUserAuth: %s", err.Error())
		return err
	}
	if user.Authenticated != 1 {
		return fmt.Errorf("without 2fa")
	}
	return nil
}

func (e EmailMessage) sendEmail() error {
	m := gomail.NewMessage()

	m.SetHeader("From", viper.GetString("email.from"))

	m.SetHeader("To", e.Email)

	m.SetHeader("Subject", e.Subject)

	m.SetBody("text/plain", fmt.Sprintf("code is %s", e.Message))

	fmt.Println(viper.GetString("email.smtp_host"), viper.GetInt("email.smtp_port"), viper.GetString("email.from"), viper.GetString("email.password"))

	d := gomail.NewDialer(viper.GetString("email.smtp_host"), viper.GetInt("email.smtp_port"), viper.GetString("email.from"), viper.GetString("email.pass"))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		slog.Info(err.Error())
		return err
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
