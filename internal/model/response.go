package model

import "github.com/golang-jwt/jwt/v5"

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
