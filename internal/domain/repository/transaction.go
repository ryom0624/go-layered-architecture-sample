package repository

import "context"

type Transaction interface {
	Commit() error
	Rollback() error
	GetDB() interface{}
}

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(tx Transaction) error) error
}