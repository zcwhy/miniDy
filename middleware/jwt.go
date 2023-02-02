package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"miniDy/model"
	"time"
)

var jwtKey = []byte("zcwhy")

type Claims struct {
	UserId int64
	jwt.RegisteredClaims
}

func ReleaseToken(user *model.UserLogin) (string, error) {
	expirationTime := time.Now().Add(15 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.UserInfoId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: expirationTime},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			Subject:   "user token",
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
