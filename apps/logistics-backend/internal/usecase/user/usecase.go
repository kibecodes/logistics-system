package user

import (
	"context"
	"fmt"
	"logistics-backend/internal/domain/driver"
	domain "logistics-backend/internal/domain/user"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repo    domain.Repository
	drvRepo domain.DriverReader
}

func NewUseCase(repo domain.Repository, drvRepo domain.DriverReader) *UseCase {
	return &UseCase{repo: repo, drvRepo: drvRepo}
}

func (uc *UseCase) RegisterUser(ctx context.Context, u *domain.User) error {
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	u.PasswordHash = string(hashedPassword)

	// insert user to DB
	if err := uc.repo.Create(u); err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}

	// if role is driver, insert into drivers table
	if u.Role == "driver" {
		driver := &driver.Driver{
			ID:              u.ID,
			FullName:        u.FullName,
			Email:           u.Email,
			VehicleInfo:     "not set",
			CurrentLocation: "not set",
			Available:       true,
			CreatedAt:       time.Now(),
		}
		if err := uc.drvRepo.RegisterDriver(ctx, driver); err != nil {
			return fmt.Errorf("could not register driver: %w", err)
		}
	}
	return nil
}

// PATCH method for drivers users to update details
func (uc *UseCase) UpdateDriverProfile(ctx context.Context, id uuid.UUID, req *domain.UpdateDriverUserProfileRequest) error {
	return uc.repo.UpdateProfile(ctx, id, req.Phone)
}

// PATCH method for user details
func (uc *UseCase) UpdateUser(ctx context.Context, userID uuid.UUID, req *domain.UpdateUserRequest) error {
	return uc.repo.UpdateColum(ctx, userID, req.Column, req.Value)
}

func (uc *UseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return uc.repo.GetByID(id)
}

func (uc *UseCase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return uc.repo.GetByEmail(email)
}

func (uc *UseCase) ListUsers(ctx context.Context) ([]*domain.User, error) {
	return uc.repo.List()
}

func (uc *UseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *UseCase) GetAllCustomers(ctx context.Context) ([]domain.AllCustomers, error) {
	return uc.repo.GetAllCustomers(ctx)
}
