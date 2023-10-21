package model

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}

type JwtCustomClaims struct {
	Login string `json:"name"`
	Id    int64  `json:"id"`
	jwt.RegisteredClaims
}

type Token struct {
	AccessToken string `json:"accessToken"`
}

type Auth struct {
	AccessToken string `json:"accessToken"`
	Id          int64  `json:"id"`
	Login       string `json:"login"`
}

type ResponseSuccess struct {
	Code   int
	Result interface{}
	Time   int64
}
type ResponseError struct {
	Code        int
	Description string
	Error       error
}

func (c *JwtCustomClaims) Valid() error {
	expirationTime := c.ExpiresAt.Time.Unix()

	if time.Now().Unix() > expirationTime {
		return errors.New("token has expired")
	}

	return nil
}
