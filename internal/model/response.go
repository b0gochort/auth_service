package model

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/restream/reindexer/v3"
)

type User struct {
	ID            int64           `json:"id"`
	Name          string          `json:"name"`
	Surname       string          `json:"surname"`
	Patronymic    string          `json:"patronymic"`
	Email         string          `json:"email"`
	Authenticated bool            `json:"authenticated"`
	Login         string          `json:"login"`
	Password      string          `json:"password"`
	IP            string          `json:"ip"`
	Birthday      int64           `json:"birthday"`
	City          string          `json:"city"`
	Position      reindexer.Point `json:"position"`
	Date          DateType        `json:"date"`
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
	Auht        bool   `json:"auth"`
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

type GeoResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Verification struct {
	Email string `db:"email"`
	Code  string `db:"code"`
}
