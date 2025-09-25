package postgres

import (
	"context"
	"fmt"
	"logistics-backend/internal/application"
	"logistics-backend/internal/domain/user"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	exec sqlx.ExtContext
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{exec: db}
}

func (r *UserRepository) execFromCtx(ctx context.Context) sqlx.ExtContext {
	if tx := application.GetTx(ctx); tx != nil {
		return tx
	}
	return r.exec
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (full_name, email, password_hash, role, status, phone, slug)
		VALUES (:full_name, :email, :password_hash, :role, :status, :phone, :slug)
		RETURNING id
	`
	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, u)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&u.ID); err != nil {
			return fmt.Errorf("scanning new user id: %w", err)
		}
	} else {
		return fmt.Errorf("no id returned after scan")
	}

	return nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, userID uuid.UUID, phone string) error {
	query := `
		UPDATE users
		SET phone = :phone, updated_at = NOW()
		WHERE id = :id
	`

	args := map[string]interface{}{
		"phone": phone,
		"id":    userID,
	}

	res, err := sqlx.NamedExecContext(ctx, r.execFromCtx(ctx), query, args)
	if err != nil {
		return fmt.Errorf("update user profile: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no user found with id %s", userID)
	}

	return nil
}

func (r *UserRepository) UpdateColum(ctx context.Context, userID uuid.UUID, column string, value any) error {
	allowed := map[string]bool{
		"full_name":  true,
		"email":      true,
		"phone":      true,
		"role":       true,
		"status":     true,
		"last_login": true,
	}

	if !allowed[column] {
		return fmt.Errorf("attempted to update disallowed column: %s", column)
	}

	query := fmt.Sprintf(`
		UPDATE users SET %s = :value, updated_at = NOW() 
		WHERE id = :id
	`, column)

	args := map[string]interface{}{
		"value": value,
		"id":    userID,
	}

	res, err := sqlx.NamedExecContext(ctx, r.execFromCtx(ctx), query, args)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no user found with id %s", userID)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, role, status, last_login, phone, slug 
		FROM users 
		WHERE id = $1
	`

	var u user.User
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &u, query, id)
	return &u, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, role, status, last_login, phone, slug 
		FROM users 
		WHERE email = $1
	`

	var u user.User
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &u, query, email)
	return &u, err
}

func (r *UserRepository) List(ctx context.Context) ([]*user.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, role, status, last_login, phone, slug 
		FROM users
	`

	var users []*user.User
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &users, query)
	return users, err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users 
		WHERE id = $1
	`
	res, err := r.exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify user deletion: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user already deleted or invalid")
	}

	return nil

}

func (r *UserRepository) GetAllCustomers(ctx context.Context) ([]user.AllCustomers, error) {
	query := `
        SELECT id, full_name
        FROM users
        WHERE role = 'customer'
        ORDER BY full_name ASC
    `
	var customers []user.AllCustomers
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &customers, query)
	return customers, err
}
