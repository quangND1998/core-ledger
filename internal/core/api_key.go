package core

import (
	config "core-ledger/configs"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

var SecretKey = config.GetConfig().Common.APISecretKey

func GenerateApiKey() (string, error) {
	const APIKeyLength = 16 // export env
	bytes := make([]byte, APIKeyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func HashApiKey(apiKey string) string {
	h := hmac.New(sha256.New, []byte(SecretKey))
	h.Write([]byte(apiKey))
	return hex.EncodeToString(h.Sum(nil))
}

func VerifyApiKey(apiKeyFromRequest, storedHash string) bool {
	hashToCheck := HashApiKey(apiKeyFromRequest)
	return hmac.Equal([]byte(hashToCheck), []byte(storedHash))
}
