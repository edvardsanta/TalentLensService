package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"platform-service/internal/config"
	"platform-service/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
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

func GenerateJWT(userID string, username string, role string, expiredAt time.Time) (string, error) {
	claims := &models.JwtCustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	secret := []byte(config.GetJWTSecretKey())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = GenerateKeyID(secret)
	return token.SignedString(secret)
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}

func JWTConfig() echojwt.Config {
	return echojwt.Config{
		SigningKey: []byte(config.GetJWTSecretKey()),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(models.JwtCustomClaims)
		},
	}
}
