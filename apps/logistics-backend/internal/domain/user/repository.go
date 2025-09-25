package user

import (
	"context"

	"github.com/google/uuid"
)

// User = actual onboarded account in the system.
type Repository interface {
	Create(ctx context.Context, user *User) error                                      // POST
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)                          // GET
	GetByEmail(ctx context.Context, email string) (*User, error)                       // GET
	List(ctx context.Context) ([]*User, error)                                         // GET
	GetAllCustomers(ctx context.Context) ([]AllCustomers, error)                       // GET
	UpdateColum(ctx context.Context, userID uuid.UUID, column string, value any) error // PATCH
	UpdateProfile(ctx context.Context, id uuid.UUID, phone string) error               // PUT
	Delete(ctx context.Context, id uuid.UUID) error                                    // DELETE
}
