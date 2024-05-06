package entity

import (
	"time"
)

type UserID int64

type User struct {
	ID             UserID    `json:"id" db:"id"`
	FirstName      string    `json:"firstName" db:"first_name"`
	FirstNameKana  string    `json:"firstNameKana" db:"first_name_kana"`
	FamilyName     string    `json:"familyName" db:"family_name"`
	FamilyNameKana string    `json:"familyName_kana" db:"family_name_kana"`
	Password       string    `json:"password" db:"password"`
	SendingPoint   int       `json:"sendingPoint" db:"sending_point"`
	Email          string    `json:"email" db:"email"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdateAt       time.Time `json:"updateAt" db:"update_at"`
}

type Users []*User
