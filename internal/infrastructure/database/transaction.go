package database

import (
	"context"
	"layered-architecture-template/internal/domain/repository"

	"gorm.io/gorm"
)

type gormTransaction struct {
	tx *gorm.DB
}

func (gt *gormTransaction) Commit() error {
	return gt.tx.Commit().Error
}

func (gt *gormTransaction) Rollback() error {
	return gt.tx.Rollback().Error
}

func (gt *gormTransaction) GetDB() interface{} {
	return gt.tx
}

type gormTransactionManager struct {
	db *gorm.DB
}

func NewGormTransactionManager(db *gorm.DB) repository.TransactionManager {
	return &gormTransactionManager{db: db}
}

func (gtm *gormTransactionManager) WithTransaction(ctx context.Context, fn func(tx repository.Transaction) error) error {
	tx := gtm.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	gormTx := &gormTransaction{tx: tx}

	defer func() {
		if r := recover(); r != nil {
			gormTx.Rollback()
			panic(r)
		}
	}()

	if err := fn(gormTx); err != nil {
		gormTx.Rollback()
		return err
	}

	return gormTx.Commit()
}