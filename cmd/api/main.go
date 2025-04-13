package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fadil-Tao/gopher-pay/config"
	"github.com/Fadil-Tao/gopher-pay/db"
	"github.com/Fadil-Tao/gopher-pay/internal/transport/rest"
	loggers "github.com/Fadil-Tao/gopher-pay/utils/logger"
)

func main() {
	slogInstance := loggers.NewLogger()
	slog.SetDefault(slogInstance)
	
	cfg := config.New() 
	Conn, err := db.InitDb(&cfg.DB)
	if err != nil {
		slog.Error("database config error", "error" , err)
	}
	defer Conn.Close() 


	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthcheck", rest.HealthCheck)
	api := http.NewServeMux()
	api.Handle("/api/", http.StripPrefix("/api", mux))
	
	server := http.Server{
		Addr:  ":" + cfg.Server.Port,
		Handler: api,
	}

	go func() {
		slog.Info("server successfully started", "Port" , server.Addr)
		if err := server.ListenAndServe(); err != nil {
			slog.Error("Server serving new connections...")
		}
		slog.Info("Stopped serving new connections")
	}()
	
	sigChan := make(chan os.Signal,1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err = server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP Shutdown error", "error", err)
	}
	slog.Info("Graceful shutdown complete")
}