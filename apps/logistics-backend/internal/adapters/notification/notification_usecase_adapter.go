package notificationadapter

import (
	notificationusecase "logistics-backend/internal/usecase/notification"
)

type UseCaseAdapter struct {
	UseCase *notificationusecase.UseCase
}
