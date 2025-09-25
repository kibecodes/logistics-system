package feedback

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, feedback *Feedback) error
	GetByID(ctx context.Context, id uuid.UUID) (*Feedback, error)
	List(ctx context.Context) ([]*Feedback, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
