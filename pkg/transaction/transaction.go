package transaction

import (
	"context"
)

type Transaction struct{}

func Begin(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func Rollback(ctx context.Context) {
}

func Commit(ctx context.Context) error {
	return nil
}

func Get(ctx context.Context) (*Transaction, error) {
	return &Transaction{}, nil
}
