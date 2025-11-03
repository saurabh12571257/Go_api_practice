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
)

func main() {
	//load config

	cfg := config.MustLoad()
	// database setup

	// setup routers
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("Welcome to Students API"))
	})

	// http server setup
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Starting server...", slog.String("addr", cfg.Addr))
	fmt.Println("Starting server on port", cfg.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}

	}()

	<-done

	slog.Info("Stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to stop server", slog.String("error", err.Error()))
	}

	slog.Info("Server stopped gracefully")

}
