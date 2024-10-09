package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"milton/pkg/db/orm"
)

type txKey struct{}

// injectTx injects transaction to context
func InjectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context
func ExtractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

func (db *SqliteDB) WithQtx(ctx context.Context) orm.Querier {
	tx := ExtractTx(ctx)
	if tx != nil {
		return orm.New(tx)
	}
	return orm.New(db.db)
}

// Transactor runs logic inside a single database transaction
type Transactor interface {
	// WithinTransaction runs a function within a database transaction.
	//
	// Transaction is propagated in the context,
	// so it is important to propagate it to underlying repositories.
	// Function commits if error is nil, and rollbacks if not.
	// It returns the same error.
	WithinTransaction(context.Context, func(ctx context.Context) error) error
}

func (db *SqliteDB) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	// begin transaction
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// run callback
	err = tFunc(InjectTx(ctx, tx))
	if err != nil {
		// if error, rollback
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("rollback transaction: %v", errRollback)
		}
		return err
	}
	// if no error, commit
	if errCommit := tx.Commit(); errCommit != nil {
		log.Printf("commit transaction: %v", errCommit)
	}
	return nil
}
