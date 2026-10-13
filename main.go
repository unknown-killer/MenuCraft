// Command ResturantMenue runs the restaurant menu service: a Go HTTP server
// with an embedded SQLite database, server-rendered bilingual templates and a
// per-venue theme engine.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ResturantMenue/internal/config"
	"ResturantMenue/internal/database"
	"ResturantMenue/internal/repository"
	"ResturantMenue/internal/server"
)

func main() {
	log.SetFlags(log.LstdFlags)

	cfg := config.Load()

	db, err := database.Open(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	store := repository.New(db)
	srv := server.New(cfg, store)

	httpServer := &http.Server{
		Addr:    cfg.Addr,
		Handler: srv.Handler(),
		// Uploads can be slow on a phone connection, so the write timeout is
		// generous while the read header timeout stays tight.
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Serve in the background so the main goroutine can wait for a signal.
	errCh := make(chan error, 1)
	go func() {
		log.Printf("menu service listening on http://localhost%s (env=%s, db=%s)", cfg.Addr, cfg.Env, cfg.DatabasePath)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Fatalf("server: %v", err)
	case sig := <-stop:
		log.Printf("received %s, shutting down", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown: %v", err)
	}
	log.Println("stopped")
}
