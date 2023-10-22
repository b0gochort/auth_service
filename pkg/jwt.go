package pkg

import (
	"time"

	"github.com/b0gochort/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(login string, id int64, expiration time.Duration) (string, error) {
	//TODO: add to config
	secretKey := []byte("salt")

	claims := &model.JwtCustomClaims{
		Login: login,
		Id:    id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
