package postgres

import (
	"context"
	"fmt"
	"logistics-backend/internal/application"
	"logistics-backend/internal/domain/store"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type StoreRepository struct {
	exec sqlx.ExtContext
}

func NewStoreRepository(db *sqlx.DB) *StoreRepository {
	return &StoreRepository{exec: db}
}

func (r *StoreRepository) execFromCtx(ctx context.Context) sqlx.ExtContext {
	if tx := application.GetTx(ctx); tx != nil {
		return tx
	}
	return r.exec
}

func (r *StoreRepository) Create(ctx context.Context, s *store.Store) error {
	query := `
        INSERT INTO stores (id, owner_id, name, slug, description, logo_url, banner_url, is_public)
        VALUES (:id, :owner_id, :name, :slug, :description, :logo_url, :banner_url, :is_public)
	`

	_, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, s)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return err
}

func (r *StoreRepository) GetByID(ctx context.Context, id uuid.UUID) (*store.Store, error) {
	query := `
		SELECT * FROM stores 
		WHERE id =  $1
	`

	var s store.Store
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &s, query, id)
	return &s, err
}

func (r *StoreRepository) GetBySlug(ctx context.Context, slug string) (*store.Store, error) {
	query := `
		SELECT * FROM stores 
		WHERE slug =  $1
	`

	var s store.Store
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &s, query, slug)
	return &s, err
}

func (r *StoreRepository) GetByOwner(ctx context.Context, ownerID uuid.UUID) (*store.Store, error) {
	query := `
		SELECT * FROM stores 
		WHERE owner_id = $1
	`

	var s store.Store
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &s, query, ownerID)
	return &s, err
}

func (r *StoreRepository) Update(ctx context.Context, storeID uuid.UUID, column string, value any) error {
	allowed := map[string]bool{
		"name":        true,
		"description": true,
		"logo_url":    true,
		"banner_url":  true,
		"slug":        true,
		"is_public":   true,
	}

	if !allowed[column] {
		return fmt.Errorf("attempted to update disallowed column: %s", column)
	}

	query := fmt.Sprintf(`
		UPDATE orders SET %s = :value, updated_at = NOW() 
		WHERE id = :id
	`, column)

	args := map[string]interface{}{
		"value": value,
		"id":    storeID,
	}

	res, err := sqlx.NamedExecContext(ctx, r.execFromCtx(ctx), query, args)
	if err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no store found with id %s", storeID)
	}

	return nil
}

func (r *StoreRepository) ListPublic(ctx context.Context) ([]*store.Store, error) {
	query := `
		SELECT * FROM stores 
		WHERE is_public = true
	`

	var stores []*store.Store
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &stores, query)
	return stores, err
}

func (r *StoreRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM stores
		WHERE id = $1
	`

	res, err := r.exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete store: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify store deletion: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("store already deleted or invalid")
	}

	return nil
}
