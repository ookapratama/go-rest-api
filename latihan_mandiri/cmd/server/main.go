package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/ookapratama/go-rest-api/latihan_mandiri/internal/api/handler"
	"github.com/ookapratama/go-rest-api/latihan_mandiri/internal/repository/postgres"
	"github.com/ookapratama/go-rest-api/latihan_mandiri/internal/service"
)

func main() {
	// 1. SETUP LOGGING (Agar log rapi dalam bentuk JSON)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 2. LOAD ENV (Membaca file .env)
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using system env")
	}

	// 3. KONEKSI DATABASE
	// Ambil URL dari .env. 
	// Untuk Supabase: Ganti DATABASE_URL di .env dengan URL dari dashboard Supabase
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		slog.Error("DATABASE_URL is not set")
		os.Exit(1)
	}

	// Buat Connection Pool
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		slog.Error("Failed to create database pool", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close() // Pastikan ditutup saat aplikasi mati

	// CARA CEK KONEKSI: Gunakan Ping!
	// Ini akan benar-benar mencoba 'salaman' dengan database.
	if err := dbPool.Ping(context.Background()); err != nil {
		slog.Error("Database UNREACHABLE", "error", err)
		os.Exit(1)
	}
	slog.Info("Successfully connected to database!")

	// 4. DEPENDENCY INJECTION (Menyambungkan Puzzle)
	// Repo -> Service -> Handler
	productRepo := postgres.NewProductRepository(dbPool)
	productSvc := service.NewProductService(productRepo)
	productHdl := handler.NewProductHandler(productSvc)

	// 5. SETUP ROUTER (CHI)
	r := chi.NewRouter()
	r.Use(middleware.Logger)    // Log setiap request masuk
	r.Use(middleware.Recoverer) // Cegah server mati jika ada error parah

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/products", func(r chi.Router) {
			r.Post("/", productHdl.Create)
			r.Get("/", productHdl.GetAll)
		})
	})

	// 6. SERVER & GRACEFUL SHUTDOWN
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Jalankan server di background (Goroutine)
	go func() {
		slog.Info("Server running", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %s\n", err)
		}
	}()

	// Menunggu sinyal STOP (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")
}
