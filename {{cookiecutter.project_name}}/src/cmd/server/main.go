package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/server"
	"honnef.co/go/tools/config"
)

const shutdownTimeout = 25 * time.Second

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := server.NewServer(addr)

	// Start the server in a separate goroutine.
	go func() {
		log.Printf("listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Block until a termination signal is received.
	<-ctx.Done()
	stop() // Allow a second Ctrl+C to force-kill.
	log.Println("shutdown signal received")

	// Mark as shutting down so readiness probes fail immediately.
	srv.SetShuttingDown()

	// Shut down gracefully, waiting for in-flight requests.
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(), shutdownTimeout,
	)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	log.Println("server stopped")
}
