package orderadapter

import (
	"context"
	"logistics-backend/internal/domain/order"
	orderusecase "logistics-backend/internal/usecase/order"

	"github.com/google/uuid"
)

type UseCaseAdapter struct {
	UseCase *orderusecase.UseCase
}

func (a *UseCaseAdapter) GetOrderByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	return a.UseCase.GetOrder(ctx, id)
}

func (a *UseCaseAdapter) UpdateOrder(ctx context.Context, orderID uuid.UUID, column string, value any) error {
	return a.UseCase.UpdateOrder(ctx, orderID, column, value)
}
