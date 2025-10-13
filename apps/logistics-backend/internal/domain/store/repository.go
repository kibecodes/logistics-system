package store

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, s *Store) error
	GetByID(ctx context.Context, id uuid.UUID) (*Store, error)
	GetBySlug(ctx context.Context, slug string) (*Store, error)
	GetByOwner(ctx context.Context, ownerID uuid.UUID) (*Store, error)
	Update(ctx context.Context, storeID uuid.UUID, column string, value any) error
	ListPublic(ctx context.Context) ([]*Store, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
