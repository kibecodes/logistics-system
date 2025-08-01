package user

import (
	generate "logistics-backend/internal/utils"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"` //raw password from client
	Role     Role   `json:"role" binding:"required,oneof=admin driver customer"`
	Phone    string `json:"phone" binding:"required"`
	Slug     string `json:"slug" binding:"required"`
}

type UpdateDriverUserProfileRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type UpdateUserRequest struct {
	Column string      `json:"column" binding:"required"`
	Value  interface{} `json:"value" binding:"required"`
}

func (r *CreateUserRequest) ToUser() *User {
	baseSlug := generate.GenerateSlug(r.FullName)
	uniqueSuffix := uuid.New().String()[:8]

	return &User{
		FullName:     r.FullName,
		Email:        r.Email,
		PasswordHash: r.Password,
		Role:         r.Role,
		Phone:        r.Phone,
		Slug:         baseSlug + "-" + uniqueSuffix,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Token    string `json:"token,omitempty"`
}
