package gormwrap

import (
	"context"

	"gorm.io/gorm"
)

// ExecuteFuncWithTxn execute a func with db txn.
func ExecuteFuncWithTxn(ctx context.Context, conn *gorm.DB, f func(tx *gorm.DB) error) (err error) {
	tx := conn.Begin().WithContext(ctx)
	if err = tx.Error; err != nil {
		return
	}

	defer func() {
		if err == nil {
			err = tx.Commit().Error
		}
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = f(tx); err != nil {
		return
	}
	return
}
