package driver

import (
	"context"
	"fmt"
	domain "logistics-backend/internal/domain/driver"
	"logistics-backend/internal/domain/notification"
	"logistics-backend/internal/usecase/common"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
)

type UseCase struct {
	repo      domain.Repository
	txManager common.TxManager
	notfRepo  domain.NotificationReader
}

func NewUseCase(repo domain.Repository, txm common.TxManager, notf domain.NotificationReader) *UseCase {
	return &UseCase{repo: repo, txManager: txm, notfRepo: notf}
}

func (uc *UseCase) RegisterDriver(ctx context.Context, d *domain.Driver) error {
	if d.ID == uuid.Nil {
		return domain.ErrMissingUserID
	}

	go func() {
		msg := fmt.Sprintf("✅ Welcome %s! Your driver account has been created.", d.FullName)
		_ = uc.notify(ctx, d.ID, msg)
	}()

	return uc.repo.Create(ctx, d)
}

func (uc *UseCase) UpdateDriverProfile(ctx context.Context, id uuid.UUID, req *domain.UpdateDriverProfileRequest) error {

	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		driver, err := uc.repo.GetByID(txCtx, id)
		if err != nil {
			return fmt.Errorf("could not fetch driver: %w", err)
		}

		if err := uc.repo.UpdateProfile(txCtx, id, req.VehicleInfo, req.CurrentLocation); err != nil {
			return fmt.Errorf("update driver profile failed: %w", err)
		}

		go func() {
			msg := "ℹ️ Your driver profile has been updated."
			_ = uc.notify(ctx, driver.ID, msg)
		}()

		return nil
	})
}

func (uc *UseCase) UpdateDriver(ctx context.Context, id uuid.UUID, req *domain.UpdateDriverRequest) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		driver, err := uc.repo.GetByID(txCtx, id)
		if err != nil {
			return fmt.Errorf("could not fetch driver: %w", err)
		}

		if err := uc.repo.UpdateColumn(txCtx, id, req.Column, req.Value); err != nil {
			return fmt.Errorf("update driver failed: %w", err)
		}

		go func() {
			msg := fmt.Sprintf("ℹ️ Your driver account column '%s' has been updated.", req.Column)
			_ = uc.notify(ctx, driver.ID, msg)
		}()

		return nil
	})
}

func (uc *UseCase) UpdateDriverAvailability(ctx context.Context, driverID uuid.UUID, column string, available bool) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		driver, err := uc.repo.GetByID(txCtx, driverID)
		if err != nil {
			return fmt.Errorf("could not fetch driver: %w", err)
		}

		if err := uc.repo.UpdateColumn(txCtx, driverID, column, available); err != nil {
			return fmt.Errorf("update driver availability failed: %w", err)
		}

		go func() {
			status := "unavailable"
			if available {
				status = "available"
			}
			msg := fmt.Sprintf("ℹ️ Your availability status has been set to '%s'.", status)
			_ = uc.notify(ctx, driver.ID, msg)
		}()

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

func (uc *UseCase) ListAvailableDrivers(ctx context.Context, available bool) ([]*domain.Driver, error) {
	return uc.repo.ListAvailableDrivers(ctx, available)
}

func (uc *UseCase) DeleteDriver(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		driver, err := uc.repo.GetByID(txCtx, id)
		if err != nil {
			return fmt.Errorf("could not fetch driver: %w", err)
		}

		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete driver failed: %w", err)
		}

		go func() {
			msg := "🗑️ Your driver account has been deleted."
			_ = uc.notify(ctx, driver.ID, msg)
		}()

		return nil
	})
}

func (uc *UseCase) GetClosestDriver(ctx context.Context, pickup postgis.PointS, maxDistance float64) (*domain.Driver, error) {
	return uc.repo.GetNearestDriver(ctx, pickup, maxDistance)
}

func (uc *UseCase) notify(ctx context.Context, userID uuid.UUID, message string) error {
	n := &notification.Notification{
		UserID:  userID,
		Message: message,
		Type:    notification.System,
		Status:  notification.Pending,
	}
	return uc.notfRepo.Create(ctx, n)
}
