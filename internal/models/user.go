package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UID          string    `gorm:"type:char(36);uniqueIndex;not null" json:"uid"`
	Username     string    `gorm:"uniqueIndex;not null;size:50"`
	Password     string    `gorm:"not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	FirstName    string    `gorm:"size:50"`
	LastName     string    `gorm:"size:50"`
	Role         string    `gorm:"default:'user';not null"`
	Status       string    `gorm:"default:'active';not null"`
	LastLogin    time.Time `gorm:"default:null"`
	LoginCount   int       `gorm:"default:0"`
	LastIP       string    `gorm:"size:45"`
	CreatedBy    uint      `gorm:"default:0"`
	UpdatedBy    uint      `gorm:"default:0"`
	DeletedBy    uint      `gorm:"default:0"`
	ProfileImage string    `gorm:"size:255"`
	CreatedAt    time.Time `gorm:"default:current_timestamp"`
	UpdatedAt    time.Time `gorm:"default:current_timestamp"`
}
type SafeUser struct {
	UID          string    `json:"uid"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	FirstName    string    `json:"firstName,omitempty"`
	LastName     string    `json:"lastName,omitempty"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	LastLogin    time.Time `json:"lastLogin"`
	LoginCount   int       `json:"loginCount"`
	ProfileImage string    `json:"profileImage,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
type JSON map[string]interface{}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.CreatedBy == 0 {
		u.CreatedBy = 1
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if u.UpdatedBy == 0 {
		u.UpdatedBy = 1
	}
	return nil
}

func (u *User) BeforeDelete(tx *gorm.DB) error {
	if u.DeletedBy == 0 {
		u.DeletedBy = 1
	}
	return nil
}

func (u *User) UpdateLastLogin(ip string) {
	u.LastLogin = time.Now()
	u.LoginCount++
	u.LastIP = ip
}

func (u *User) IsActive() bool {
	return u.Status == "active"
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

func (u *User) SetStatus(status string) {
	u.Status = status
}

func (u *User) ToSafeUser() SafeUser {
	return SafeUser{
		UID:          u.UID,
		Username:     u.Username,
		Email:        u.Email,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Role:         u.Role,
		Status:       u.Status,
		LastLogin:    u.LastLogin,
		LoginCount:   u.LoginCount,
		ProfileImage: u.ProfileImage,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
