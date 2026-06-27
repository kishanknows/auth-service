package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateRefreshToken() (string, string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)

	if err != nil {
		return "", "", err
	}

	refreshToken := hex.EncodeToString(token)
	tokenHash := GenerateTokenHash(refreshToken)

	return refreshToken, tokenHash, nil
}

func VerifyRefreshToken(refreshToken string, tokenHash string) (bool, error) {
	refreshTokenHash := GenerateTokenHash(refreshToken)

	if refreshTokenHash == tokenHash {
		return true, nil
	}

	return false, nil
}

func GenerateTokenHash(token string) string {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])
	return tokenHash
}