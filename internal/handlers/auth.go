package handlers

import (
	"net/http"
	"platform-service/internal/database"
	"platform-service/internal/models"
	"platform-service/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username     string `json:"username" validate:"required,min=3,max=50"`
	Password     string `json:"password" validate:"required,min=6"`
	Email        string `json:"email" validate:"required,email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	ProfileImage string `json:"profile_image"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserInfo  `json:"user"`
}

type UserInfo struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	LastLogin time.Time `json:"last_login"`
}

func Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	var existingUser models.User
	if err := database.DB.Where("LOWER(username) = LOWER(?) OR LOWER(email) = LOWER(?)",
		req.Username, req.Email).First(&existingUser).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "Username or email already exists",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to hash password",
		})
	}

	user := &models.User{
		Username:     req.Username,
		Password:     string(hashedPassword),
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		ProfileImage: req.ProfileImage,
		Role:         "user",
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastIP:       c.RealIP(),
		UID:          uuid.NewString(),
	}

	tx := database.DB.Begin()
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user",
		})
	}

	if err := tx.Commit().Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to commit transaction",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
		"userId":  user.UID,
	})
}

func Login(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	storedUser := new(models.User)
	result := database.DB.Where("username = ?", user.Username).First(storedUser)
	if result.Error == gorm.ErrRecordNotFound {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	} else if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user"})
	}

	if !storedUser.IsActive() {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Account is not active"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}
	expiredAt := time.Now().Add(time.Hour * 24)
	token, err := utils.GenerateJWT(storedUser.UID, storedUser.Username, storedUser.Role, expiredAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	storedUser.UpdateLastLogin(c.RealIP())
	database.DB.Save(storedUser)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  storedUser.ToSafeUser(),
	})
}
