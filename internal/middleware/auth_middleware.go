package internal_middleware

import (
	"net/http"
	"platform-service/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*models.JwtCustomClaims)

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		return next(c)
	}
}

func AdminAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*models.JwtCustomClaims)
		if claims.Role != "admin" {
			return echo.NewHTTPError(http.StatusForbidden, "Admin access required")
		}
		return next(c)
	}
}
