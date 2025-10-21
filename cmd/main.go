package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"koois_core/internal/config"
	"koois_core/internal/db"
	"koois_core/internal/handler"
	"koois_core/internal/middleware"
	"koois_core/internal/service"
)

func main() {
	// Load and validate config
	cfg := config.Load()
	cfg.Validate()

	// Initialize database
	pool, err := db.NewPool(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer pool.Close()

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := pool.Ping(ctx); err != nil {
		cancel()
		log.Fatalf("Failed to ping database: %v", err)
	}
	cancel()

	// Initialize services
	userService := service.NewUserService(pool)
	quizService := service.NewQuizService(pool)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	quizHandler := handler.NewQuizHandler(quizService)
	fileHandler := handler.FileHandler{}

	// Setup router
	router := http.NewServeMux()

	// Public routes
	router.HandleFunc("GET /api/users", middleware.JWT(userHandler.GetAll, cfg))
	router.HandleFunc("POST /api/upload", middleware.JWT(fileHandler.Create, cfg))
	router.HandleFunc("DELETE /api/upload/{id}", middleware.JWT(fileHandler.Delete, cfg))
	router.HandleFunc("GET /api/quiz/{id}", middleware.JWT(quizHandler.GetByID, cfg))
	router.HandleFunc("POST /api/quiz", middleware.JWT(quizHandler.Create, cfg))

	// Health check
	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// 404 handler for undefined routes
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"route not found"}`))
	})

	// Server setup
	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
		pool.Close()
	}()

	log.Printf("Server starting on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
