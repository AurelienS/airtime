package storage

import (
	"context"

	"github.com/AurelienS/cigare/internal/log"
	"github.com/jackc/pgx/v5"
)

type TransactionManager struct {
	db *pgx.Conn
}

func NewTransactionManager(db *pgx.Conn) TransactionManager {
	return TransactionManager{db: db}
}

func (tm *TransactionManager) ExecuteTransaction(ctx context.Context, txFunc func() error) error {
	tx, err := tm.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Error().Msgf("rollback error: %v", rbErr)
			}
		}
	}()

	if err = txFunc(); err != nil {
		log.Error().Msgf("transaction error: %v", err)
		return err
	}

	return tx.Commit(ctx)
}
