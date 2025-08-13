package orderadapter

import (
	"context"
	"logistics-backend/internal/domain/order"
	orderusecase "logistics-backend/internal/usecase/order"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UseCaseAdapter struct {
	UseCase *orderusecase.UseCase
}

func (a *UseCaseAdapter) GetOrderByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	return a.UseCase.GetOrder(ctx, id)
}

func (a *UseCaseAdapter) UpdateOrderTx(ctx context.Context, tx *sqlx.Tx, orderID uuid.UUID, column string, value any) error {
	return a.UseCase.UpdateOrderTx(ctx, tx, orderID, column, value)
}
