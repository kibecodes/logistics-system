package driver

import (
	"context"
	"fmt"
	"log"
	domain "logistics-backend/internal/domain/driver"

	"github.com/google/uuid"
)

type UseCase struct {
	repo domain.Repository
}

func NewUseCase(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

// -- might not actually need this.
func (uc *UseCase) RegisterDriver(ctx context.Context, d *domain.Driver) error {
	// Validate required fields
	if d.ID == uuid.Nil {
		return domain.ErrMissingUserID
	}

	return uc.repo.Create(ctx, d)
}

func (uc *UseCase) UpdateDriverProfile(ctx context.Context, id uuid.UUID, req *domain.UpdateDriverProfileRequest) error {
	driver, err := uc.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get driver by ID: %w", err)
	}
	if driver == nil {
		return fmt.Errorf("driver not found with ID: %s", id)
	}

	return uc.repo.UpdateProfile(ctx, id, req.VehicleInfo, req.CurrentLocation)
}

func (uc *UseCase) UpdateDriver(ctx context.Context, id uuid.UUID, req *domain.UpdateDriverRequest) error {
	return uc.repo.UpdateColumn(ctx, id, req.Column, req.Value)
}

func (uc *UseCase) UpdateDriverAvailability(ctx context.Context, driverID uuid.UUID, column string, available bool) error {
	return uc.repo.UpdateColumn(ctx, driverID, column, available)
}

func (uc *UseCase) GetDriver(ctx context.Context, id uuid.UUID) (*domain.Driver, error) {
	log.Printf("driver usecase driver id: %+v", id)
	driver, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	log.Printf("driver from db: %+v", driver)
	return driver, nil
}

func (uc *UseCase) GetDriverByEmail(ctx context.Context, email string) (*domain.Driver, error) {
	return uc.repo.GetByEmail(email)
}

func (uc *UseCase) ListDrivers(ctx context.Context) ([]*domain.Driver, error) {
	return uc.repo.List()
}

func (uc *UseCase) DeleteDriver(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}
