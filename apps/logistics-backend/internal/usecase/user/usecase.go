package user

import (
	"context"
	"fmt"
	"logistics-backend/internal/domain/driver"
	domain "logistics-backend/internal/domain/user"
	"logistics-backend/internal/usecase/common"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repo      domain.Repository
	drvRepo   domain.DriverReader
	txManager common.TxManager
}

func NewUseCase(repo domain.Repository, drvRepo domain.DriverReader, txm common.TxManager) *UseCase {
	return &UseCase{repo: repo, drvRepo: drvRepo, txManager: txm}
}

func (uc *UseCase) RegisterUser(ctx context.Context, u *domain.User) error {

	return uc.txManager.Do(ctx, func(txCtx context.Context) error {

		// 1. hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		u.PasswordHash = string(hashedPassword)

		// 2. insert user to DB
		if err := uc.repo.Create(txCtx, u); err != nil {
			return fmt.Errorf("could not create user: %w", err)
		}

		// 3. if role is driver, insert into drivers table
		if u.Role == "driver" {
			driver := &driver.Driver{
				ID:          u.ID,
				FullName:    u.FullName,
				Email:       u.Email,
				VehicleInfo: "not set",
				CurrentLocation: postgis.PointS{
					SRID: 4326,
					X:    36.8219, // longitude
					Y:    -1.2921, // latitude
				},
				Available: true,
				CreatedAt: time.Now(),
			}
			if err := uc.drvRepo.RegisterDriver(txCtx, driver); err != nil {
				return fmt.Errorf("could not register driver: %w", err)
			}
		}
		return nil
	})
}

// PATCH method for users to update details
func (uc *UseCase) UpdateUserProfile(ctx context.Context, id uuid.UUID, req *domain.UpdateDriverUserProfileRequest) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.UpdateProfile(txCtx, id, req.Phone); err != nil {
			return fmt.Errorf("update user profile failed: %w", err)
		}

		return nil
	})
}

// PATCH method for user details
func (uc *UseCase) UpdateUser(ctx context.Context, userID uuid.UUID, req *domain.UpdateUserRequest) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.UpdateColum(txCtx, userID, req.Column, req.Value); err != nil {
			return fmt.Errorf("update user failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return uc.repo.GetByEmail(ctx, email)
}

func (uc *UseCase) ListUsers(ctx context.Context) ([]*domain.User, error) {
	return uc.repo.List(ctx)
}

func (uc *UseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete user failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) GetAllCustomers(ctx context.Context) ([]domain.AllCustomers, error) {
	return uc.repo.GetAllCustomers(ctx)
}
