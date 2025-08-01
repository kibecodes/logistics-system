package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(user *User) error                                                           // POST
	GetByID(id uuid.UUID) (*User, error)                                               // GET
	GetByEmail(email string) (*User, error)                                            // GET
	List() ([]*User, error)                                                            // GET
	UpdateColum(ctx context.Context, userID uuid.UUID, column string, value any) error // PATCH                                         // PATCH
	UpdateProfile(ctx context.Context, id uuid.UUID, phone string) error               // PUT
	Delete(ctx context.Context, id uuid.UUID) error                                    // DELETE
}
