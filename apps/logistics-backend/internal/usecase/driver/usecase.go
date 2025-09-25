package driver

import (
	"context"
	"fmt"
	domain "logistics-backend/internal/domain/driver"
	"logistics-backend/internal/usecase/common"

	"github.com/google/uuid"
)

type UseCase struct {
	repo      domain.Repository
	txManager common.TxManager
}

func NewUseCase(repo domain.Repository, txm common.TxManager) *UseCase {
	return &UseCase{repo: repo, txManager: txm}
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

	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.UpdateProfile(txCtx, id, req.VehicleInfo, req.CurrentLocation); err != nil {
			return fmt.Errorf("update driver profile failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) UpdateDriver(ctx context.Context, id uuid.UUID, req *domain.UpdateDriverRequest) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.UpdateColumn(txCtx, id, req.Column, req.Value); err != nil {
			return fmt.Errorf("update driver failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) UpdateDriverAvailability(ctx context.Context, driverID uuid.UUID, column string, available bool) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.UpdateColumn(txCtx, driverID, column, available); err != nil {
			return fmt.Errorf("update driver availability failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) GetDriver(ctx context.Context, id uuid.UUID) (*domain.Driver, error) {
	driver, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func (uc *UseCase) GetDriverByEmail(ctx context.Context, email string) (*domain.Driver, error) {
	return uc.repo.GetByEmail(ctx, email)
}

func (uc *UseCase) ListDrivers(ctx context.Context) ([]*domain.Driver, error) {
	return uc.repo.List(ctx)
}

func (uc *UseCase) DeleteDriver(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete driver failed: %w", err)
		}

		return nil
	})
}
