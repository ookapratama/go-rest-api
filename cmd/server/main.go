package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/ookapratama/go-rest-api/internal/api"
	"github.com/ookapratama/go-rest-api/internal/api/handler"
	"github.com/ookapratama/go-rest-api/internal/repository/postgres"
	"github.com/ookapratama/go-rest-api/internal/service"
)

func main() {
	// 1. Setup Logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 2. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using defaults")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		slog.Error("DATABASE_URL must be set")
		os.Exit(1)
	}

	// 3. Connect to Database
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Ping to verify connection
	if err := dbPool.Ping(context.Background()); err != nil {
		slog.Error("Database ping failed", "error", err)
		os.Exit(1)
	}
	slog.Info("Successfully connected to database")

	// 4. Initialize Dependency Injection
	productRepo := postgres.NewProductRepository(dbPool)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	// 5. Setup Router
	router := api.NewRouter(productHandler)

	// 6. Start Server with Graceful Shutdown
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Channel to listen for errors or interrupt signals
	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("Server is starting", "port", port)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking select
	select {
	case err := <-serverErrors:
		slog.Error("Server error", "error", err)
	case sig := <-shutdown:
		slog.Info("Shutdown signal received", "signal", sig)

		// Create context for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Graceful shutdown failed", "error", err)
			server.Close()
		}
	}

	slog.Info("Server stopped")
}
