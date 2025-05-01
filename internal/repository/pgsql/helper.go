package pgsql

import (
	"database/sql"
	"log/slog"
)


func runInTx(db *sql.DB, fn func(tx *sql.Tx)error) error{
	tx, err := db.Begin()
	if err != nil {
		slog.Error(err.Error())
		return err
	}	

	err = fn(tx)
	if err == nil {
		return tx.Commit()
	}	

	rollBackErr := tx.Rollback()
	if rollBackErr != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}