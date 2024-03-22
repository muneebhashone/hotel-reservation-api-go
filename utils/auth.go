package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/muneebhashone/go-fiber-api/config"
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token
	t, err := token.SignedString(config.JWT_SECRET)
	return t, err
}
