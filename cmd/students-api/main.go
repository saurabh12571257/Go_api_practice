package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saurabh/students-api/internal/config"
	student "github.com/saurabh/students-api/internal/http/handlers/students"
	"github.com/saurabh/students-api/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	slog.Info("strorage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	router := http.NewServeMux()
	router.HandleFunc("/api/students", student.New(storage))

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Starting server", "addr", cfg.Addr)
	fmt.Println("Starting server on port", cfg.Addr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-done
	slog.Info("Stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to stop server", "error", err)
	}

	slog.Info("Server stopped gracefully")
}
