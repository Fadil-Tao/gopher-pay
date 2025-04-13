package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/Fadil-Tao/gopher-pay/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	migrationsPath := flag.String("path", "db/migrations", "Path to migrations directory")
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	flag.Parse()

	c := config.NewDB()

	dsn := "postgres://"+c.Username+ ":" + c.Password + "@" + c.Host + ":" + strconv.Itoa(c.Port) + "/" + c.DBName + "?sslmode=disable"

	slog.Info("Running migrations", "command", command, "dsn", dsn)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", *migrationsPath),
		dsn,
	)
	if err != nil {
		slog.Error("Failed to initialize migrate", "err", err)
		os.Exit(1)
	}

	defer m.Close()

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			slog.Error("Migration up failed", "err", err)
			os.Exit(1)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			slog.Error("Migration down failed", "err", err)
			os.Exit(1)
		}
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			slog.Error("Failed to get version", "err", err)
			os.Exit(1)
		}
		fmt.Printf("Version: %d, Dirty: %t\n", version, dirty)
	default:
		slog.Error("Unknown command", "command", command)
		os.Exit(1)
	}

	slog.Info("Migration completed successfully", "command", command)
}