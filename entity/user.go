package entity

import (
	"time"
)

type UserID int64

type User struct {
	ID         UserID    `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Password   string    `json:"password" db:"password"`
	Email      string    `json:"email" db:"email"`
	Role       int       `json:"role" db:"role_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	ModifiedAt time.Time `json:"modifiedAt" db:"modified_at"`
}
