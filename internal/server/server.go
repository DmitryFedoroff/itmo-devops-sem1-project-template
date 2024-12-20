package server

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

	"project_sem/internal/api"
	"project_sem/platform/config"
	"project_sem/platform/storage"
)

type Application interface {
	Run()
}

type application struct {
	server *http.Server
}

func New(cfg config.Settings) Application {
	repo, err := storage.NewRepository(cfg.Database)
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/prices", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			api.PostPrices(repo)(w, r)
			return
		} else if r.Method == http.MethodGet {
			api.GetPrices(repo)(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	serverInstance := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &application{server: serverInstance}
}

func (a *application) Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Printf("starting server on %s", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down gracefully...")

	const defaultShutdownTimeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	log.Println("server shutdown complete")
}
