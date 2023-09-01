package util

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func CleanupTx(tx *sql.Tx) {
	if p := recover(); p != nil { //nolint:gocritic
		tx.Rollback() //nolint:errcheck
		panic(p)      // re-throw panic after Rollback
	}

	// if a transaction has already been commited
	// it will be a no op
	tx.Rollback() //nolint:errcheck
}

// https://stackoverflow.com/questions/16184238/database-sql-tx-detecting-commit-or-rollback
func Transact(db *sqlx.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil { //nolint:gocritic
			tx.Rollback() //nolint:errcheck
			panic(p)      // re-throw panic after Rollback
		} else if err != nil {
			// err is non-nil; don't change it
			tx.Rollback() // nolint:errcheck
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()

	err = txFunc(tx)

	return err
}
