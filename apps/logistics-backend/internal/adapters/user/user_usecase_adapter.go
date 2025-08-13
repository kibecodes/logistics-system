package useradapter

import (
	"context"
	"logistics-backend/internal/domain/order"
	userusecase "logistics-backend/internal/usecase/user"
)

type UseCaseAdapter struct {
	UseCase *userusecase.UseCase
}

func (a *UseCaseAdapter) GetAllCustomers(ctx context.Context) ([]order.Customer, error) {
	custs, err := a.UseCase.GetAllCustomers(ctx) // returns []user.AllCustomers
	if err != nil {
		return nil, err
	}

	// Map []user.AllCustomers → []order.Customer
	res := make([]order.Customer, len(custs))
	for i, c := range custs {
		res[i] = order.Customer{
			ID:   c.ID,
			Name: c.Name,
			// map only what Order domain needs
		}
	}
	return res, nil
}
