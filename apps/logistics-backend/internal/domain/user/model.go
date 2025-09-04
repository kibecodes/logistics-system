package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Role string
type UserStatus string

const (
	Admin    Role = "admin"
	Driver   Role = "driver"
	Customer Role = "customer"
	Guest    Role = "guest"
)

const (
	Active    UserStatus = "active"
	Inactive  UserStatus = "inactive"
	Suspended UserStatus = "suspended"
	Pending   UserStatus = "pending"
)

type User struct {
	ID                   uuid.UUID  `db:"id" json:"id"`
	FullName             string     `db:"full_name" json:"fullName"`
	Email                string     `db:"email" json:"email"`
	PasswordHash         string     `db:"password_hash" json:"-"`
	Role                 Role       `db:"role" json:"role"`
	Phone                string     `db:"phone" json:"phone"`
	Slug                 string     `db:"slug" json:"slug"` // adminSlug used in public route
	Must_change_password bool       `db:"must_change_password" json:"must_change_password"`
	Status               UserStatus `db:"status" json:"status"`
	LastLogin            *time.Time `db:"last_login" json:"last_login,omitempty"`
	CreatedAt            time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at" json:"updated_at"`
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

type AllCustomers struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"full_name" json:"name"`
}
