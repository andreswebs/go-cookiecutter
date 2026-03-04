// Package server provides HTTP server setup with graceful shutdown support.
package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync/atomic"
)

// Server wraps an [http.Server] with graceful shutdown support.
// It tracks shutdown state so that readiness probes can fail early
// once a termination signal is received.
type Server struct {
	httpServer     *http.Server
	isShuttingDown atomic.Bool
}

// NewServer creates a new Server bound to the given address.
func NewServer(addr string) *Server {
	s := &Server{}
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.routes(),
	}
	return s
}

// Handler returns the server's HTTP handler for use in tests.
func (s *Server) Handler() http.Handler {
	return s.httpServer.Handler
}

// ListenAndServe starts accepting connections. It blocks until the server
// is shut down and returns [http.ErrServerClosed] on graceful shutdown.
func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server. It marks the server as
// shutting down (causing readiness probes to fail), then waits for
// in-flight requests to complete or the context to expire.
func (s *Server) Shutdown(ctx context.Context) error {
	s.isShuttingDown.Store(true)
	return s.httpServer.Shutdown(ctx)
}

// SetShuttingDown marks the server as shutting down. Subsequent health
// checks will return 503 Service Unavailable. This is useful for
// signaling load balancers before calling [Server.Shutdown].
func (s *Server) SetShuttingDown() {
	s.isShuttingDown.Store(true)
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /livez", s.handleLiveness)
	mux.HandleFunc("GET /readyz", s.handleReadiness)
	return mux
}

// handleLiveness returns 200 OK if the process can handle HTTP requests.
// No dependency checks — fast and side-effect-free.
func (s *Server) handleLiveness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		slog.Error("failed to encode liveness response", "error", err)
	}
}

// handleReadiness returns 200 OK if the service can serve requests.
// Returns 503 Service Unavailable during shutdown.
func (s *Server) handleReadiness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if s.isShuttingDown.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "unavailable"}); err != nil {
			slog.Error("failed to encode readiness response", "error", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		slog.Error("failed to encode readiness response", "error", err)
	}
}
