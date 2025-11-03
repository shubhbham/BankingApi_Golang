package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shubhbham/BankingApi_Golang/internal/api"
	"github.com/shubhbham/BankingApi_Golang/internal/config"
	"github.com/shubhbham/BankingApi_Golang/internal/db"
)

type Server struct {
	config *config.Config
	db     *db.DB
	http   *http.Server
}

func New(cfg *config.Config) (*Server, error) {
	// Initialize database
	database, err := db.New(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize router
	router := api.NewRouter(database)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		config: cfg,
		db:     database,
		http:   httpServer,
	}, nil
}

func (s *Server) Start() error {
	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", s.config.ServerPort)
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.http.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return err
	}

	s.db.Close()
	log.Println("Server stopped")

	return nil
}