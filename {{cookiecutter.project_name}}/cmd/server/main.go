// Command server runs the {{ cookiecutter.project_name }} HTTP API.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/config"
	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/server"
	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/version"
)

const shutdownTimeout = 25 * time.Second

// Exit codes: a long-running server adopts the applicable slice of the exit
// taxonomy only. 78 (EX_CONFIG) for a configuration failure at startup, 1 for
// a runtime server failure, and 128 plus the signal number (130 SIGINT, 143
// SIGTERM) after a caught signal and graceful shutdown.
const exitConfig = 78

// main is the single exit boundary: run returns the exit code and every
// deferred cleanup inside it has already executed.
func main() {
	os.Exit(run())
}

func run() int {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("failed to load config: %v", err)
		return exitConfig
	}

	log.Printf("starting {{ cookiecutter.project_name }} %s", version.Current())

	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := server.NewServer(addr)

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// Catch termination signals explicitly (not via signal.NotifyContext) so
	// the exit code can report which signal ended the process.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	var exit int
	select {
	case err := <-serverErr:
		log.Printf("server error: %v", err)
		return 1
	case sig := <-sigCh:
		signal.Stop(sigCh) // Allow a second Ctrl+C to force-kill.
		log.Printf("shutdown signal received: %v", sig)
		// A tool that catches SIGINT/SIGTERM for graceful shutdown exits 128
		// plus the signal number after cleanup.
		exit = 128 + int(sig.(syscall.Signal))
	}

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
	return exit
}
