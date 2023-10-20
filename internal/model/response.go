package model

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}

type JwtCustomClaims struct {
	Name         string `json:"name"`
	RefreshToken string `json:"refresh_token"`
	jwt.RegisteredClaims
}

type Tokens struct {
	AccessToken string `json:"accessToken"`
}
