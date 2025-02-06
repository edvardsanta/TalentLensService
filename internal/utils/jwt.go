package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"platform-service/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func init() {
	jwtSecret = []byte(config.GetJWTSecretKey())
	if len(jwtSecret) == 0 {
		panic("JWT_SECRET_KEY is not set!")
	}
}

func GenerateKeyID(secret []byte) string {
	hash := sha256.Sum256(secret)
	return hex.EncodeToString(hash[:8])
}

func GenerateJWT(userID uint, username string, role string, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
	})
	token.Header["kid"] = GenerateKeyID(jwtSecret)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}
