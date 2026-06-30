package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/ArtemPotapov52/gpurenta-backend/internal/config"
	"github.com/ArtemPotapov52/gpurenta-backend/internal/db"
	"github.com/ArtemPotapov52/gpurenta-backend/internal/handler"
	"github.com/ArtemPotapov52/gpurenta-backend/internal/middleware"
)

func main() {
	cfg := config.Load()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
	slog.Info("starting server", "port", cfg.Port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	store, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	slog.Info("connected to database")

	if err := store.RunMigrations(ctx); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("migrations applied")

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := store.MarkStaleAgentsOffline(context.Background(), 5*time.Minute); err != nil {
				slog.Error("failed to mark stale agents offline", "error", err)
			}
		}
	}()

	authH := &handler.AuthHandler{Store: store, JWTSecret: cfg.JWTSecret}
	agentH := &handler.AgentHandler{Store: store}
	gpuH := &handler.GPUHandler{Store: store}
	rentalH := &handler.RentalHandler{Store: store}

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.CORS(cfg.FrontendURL))
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Timeout(30 * time.Second))

	r.Get("/v1/health", handler.Health)
	r.Get("/v1/images", agentH.GetSupportedImages)

	r.Post("/v1/auth/google", authH.Google)
	r.Post("/v1/auth/dev", authH.DevLogin)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		r.Post("/v1/agents/register", agentH.Register)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AgentAuthMiddleware)
		r.Post("/v1/agents/heartbeat", agentH.Heartbeat)
	})

	r.Get("/v1/tokens/validate", rentalH.ValidateToken)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		r.Get("/v1/gpus", gpuH.List)
		r.Post("/v1/rentals/start", rentalH.Start)
		r.Post("/v1/rentals/{id}/stop", rentalH.Stop)
		r.Get("/v1/rentals/{id}", rentalH.Get)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AgentAuthMiddleware)
		r.Get("/v1/agents/{id}/rentals", rentalH.ListByAgent)
	})

	webDir := filepath.Join("web")
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("forced shutdown", "error", err)
	}
	slog.Info("server stopped")
}
