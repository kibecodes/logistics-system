package common

import "context"

// TxManager is an abstraction usecases depend on.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
