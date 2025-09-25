package application

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type SQLTxManager struct {
	db *sqlx.DB
}

// Unexported key type for context safety
type txCtxKey struct{}

func NewTxManager(db *sqlx.DB) *SQLTxManager {
	return &SQLTxManager{db: db}
}

func (m *SQLTxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// Store the transaction in a new ctx
	txCtx := context.WithValue(ctx, txCtxKey{}, tx)

	// Run the business logic
	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("rollback failed: %v", rbErr)
		}
		return err
	}

	return tx.Commit()
}

func GetTx(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(txCtxKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}
