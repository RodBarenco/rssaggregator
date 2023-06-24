package database

import (
	"crypto/rand"
	"encoding/hex"

	"gorm.io/gorm"
)

func GetUserByAPIKey(db *gorm.DB, apiKey string) (User, error) {
	var user User
	if err := db.Where("api_key = ?", apiKey).First(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func generateAPIKey() (string, error) {
	apiKeyBytes := make([]byte, 32)
	if _, err := rand.Read(apiKeyBytes); err != nil {
		return "", err
	}

	apiKey := hex.EncodeToString(apiKeyBytes)
	return apiKey, nil
}
