package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	Admin    Role = "admin"
	Driver   Role = "driver"
	Customer Role = "customer"
)

type User struct {
	ID                   uuid.UUID `db:"id" json:"id"`
	FullName             string    `db:"full_name" json:"fullName"`
	Email                string    `db:"email" json:"email"`
	PasswordHash         string    `db:"password_hash" json:"-"`
	Role                 Role      `db:"role" json:"role"`
	Phone                string    `db:"phone" json:"phone"`
	Slug                 string    `db:"slug" json:"slug"` // adminSlug used in public route
	Must_change_password bool      `db:"must_change_password" json:"must_change_password"`
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time `db:"updated_at" json:"updated_at"`
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

type AllCustomers struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"full_name" json:"name"`
}
