package models

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	ExpiresAt time.Time
}

func GenerateToken(db *gorm.DB) (string, error) {
	// token is a random string of 8 characters
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	const tokenLength = 8

	b := make([]byte, tokenLength)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(time.Nanosecond) // ensure different seed for each iteration
	}
	tokenStr := string(b)

	expiresAt := time.Now().Add(8 * time.Hour)

	token := Token{
		Token:     tokenStr,
		ExpiresAt: expiresAt,
	}

	err := db.Create(&token).Error
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func ValidateToken(db *gorm.DB, token string) (bool, error) {
	var tokenModel Token
	err := db.Where("token = ?", token).First(&tokenModel).Error
	if err != nil {
		return false, err
	}

	if tokenModel.ExpiresAt.Before(time.Now()) {
		db.Delete(&tokenModel)
		return false, nil
	}

	return true, nil
}
