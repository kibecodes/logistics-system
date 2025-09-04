package postgres

import (
	"context"
	"fmt"
	"logistics-backend/internal/domain/user"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) user.Repository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(u *user.User) error {
	query := `
		INSERT INTO users (full_name, email, password_hash, role, status, phone, slug)
		VALUES (:full_name, :email, :password_hash, :role, :status, :phone, :slug)
		RETURNING id
	`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&u.ID, u)
}

func (r *UserRepository) UpdateProfile(ctx context.Context, id uuid.UUID, phone string) error {
	query := `
		UPDATE users
		SET phone = $1, updated_at = NOW()
		WHERE id = $2
	`

	res, err := r.db.ExecContext(ctx, query, phone, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no user found with id %s", id)
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

	query := fmt.Sprintf(`UPDATE users SET %s = $1, updated_at = NOW() WHERE id = $2`, column)
	res, err := r.db.ExecContext(ctx, query, value, userID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no user found with id %s", userID)
	}
	return nil
}

func (r *UserRepository) GetByID(id uuid.UUID) (*user.User, error) {
	query := `SELECT id, full_name, email, password_hash, role, status, last_login, phone, slug FROM users WHERE id = $1`
	var u user.User
	err := r.db.Get(&u, query, id)
	return &u, err
}

func (r *UserRepository) GetByEmail(email string) (*user.User, error) {
	query := `SELECT id, full_name, email, password_hash, role, status, last_login, phone, slug FROM users WHERE email = $1`
	var u user.User
	err := r.db.Get(&u, query, email)
	return &u, err
}

func (r *UserRepository) List() ([]*user.User, error) {
	query := `SELECT id, full_name, email, password_hash, role, status, last_login, phone, slug FROM users`
	var users []*user.User
	err := r.db.Select(&users, query)
	return users, err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepository) GetAllCustomers(ctx context.Context) ([]user.AllCustomers, error) {
	query := `
        SELECT id, full_name
        FROM users
        WHERE role = 'customer'
        ORDER BY full_name ASC
    `
	var customers []user.AllCustomers
	err := r.db.Select(&customers, query)
	return customers, err
}
