package utils

import (
	"crypto/sha512"
	"fmt"
)

func GenerateSessionToken(username, secret, password string) string {
	data := fmt.Sprintf("%s:%s:%s", username, secret, password)
	hash := sha512.New().Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func ValidateSession(token, username, secret, password string) bool {
	expectedToken := GenerateSessionToken(username, secret, password)
	return token == expectedToken
}
