package main

import (
	"context"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/config"
	httptransport "github.com/gatiyatullin-ibragim/auto-recruiter/internal/transport/http"
	"github.com/gatiyatullin-ibragim/auto-recruiter/pkg/logger"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	appLogger := logger.New(cfg.Environment)
	server := &http.Server{Addr: cfg.HTTPAddress, Handler: httptransport.NewServer(appLogger).Handler()}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			appLogger.Error("server shutdown failed", "error", err)
		}
	}()
	appLogger.Info("starting server", "addr", cfg.HTTPAddress)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		appLogger.Error("server stopped", "error", err)
	}
}
