package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(source string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func Sign(secretID, secretKey, iss, aud string) string {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Add(0).Unix(),
		"aud": aud,
		"iss": iss,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = secretID

	tokenString, _ := token.SignedString([]byte(secretKey))

	return tokenString
}
