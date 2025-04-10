package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/Fadil-Tao/gopher-pay/config"
)

func InitDb(cfg *config.ConfDB) (*sql.DB, error) {
	dsn := "postgres://" + cfg.Username + ":" + cfg.Password + "@" + cfg.Host + ":" + strconv.Itoa(cfg.Port) + "/" + cfg.DBName
	db, err := sql.Open("libsql", dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", dsn, err)
		return nil, err
	}
	slog.Info("Succesfully connect to database")

	db.SetConnMaxLifetime(cfg.MaxIdleTime * time.Minute)
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetMaxOpenConns(cfg.MaxOpenCOnn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
