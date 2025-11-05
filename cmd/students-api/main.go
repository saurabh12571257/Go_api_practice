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
)

func main() {
	cfg := config.MustLoad()

	router := http.NewServeMux()
	router.HandleFunc("/api/students", student.New())

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
