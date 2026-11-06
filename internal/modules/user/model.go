package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:100;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;size:150;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Role      string         `gorm:"size:50;default:'user'" json:"role"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email,max=150"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type UpdateUserRequest struct {
	Username *string `json:"username" binding:"omitempty,min=3,max=100"`
	Email    *string `json:"email" binding:"omitempty,email,max=150"`
	Avatar   *string `json:"avatar" binding:"omitempty,max=255"`
}
