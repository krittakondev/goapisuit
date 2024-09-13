package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	Sub         int      `json:"sub"`
	Permissions []string `json:"permissions"`
}

func SignJwt(jwt_claim JwtClaims) (token_string string, err error) {
	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		err = errors.New("JWT_SECRET not found")
	}
	jwt_expire := "1d"
	if expire := os.Getenv("JWT_EXPIRE"); expire != "" {
		jwt_secret = expire
	}
	duration, err := time.ParseDuration(jwt_expire)
	if err != nil {
		return
	}

	expirationTime := time.Now().Add(duration) // กำหนดเวลาหมดอายุ 5 นาที

	jwt_claim.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt_claim)
	token_string, err = token.SignedString([]byte(jwt_secret))

	return
}
func GenerateSecret(length int) (string, error) {
	secret := make([]byte, length)

	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(secret), nil
}
