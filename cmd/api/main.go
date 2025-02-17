package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"WB_GO_L0/internal/database"
	"WB_GO_L0/internal/kafka"
	"WB_GO_L0/internal/server"

	"golang.org/x/sync/errgroup"
)

func gracefulShutdown(apiServer *http.Server, cancelKafka context.CancelFunc) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		slog.Info("Server forced to shutdown with", "error", err)
	}

	cancelKafka()
	slog.Info("Server exiting")

	return nil
}

func main() {
	setupLogger()
	slog.Info("Starting service...")
	server := server.NewServer()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eg := new(errgroup.Group)

	eg.Go(func() error { return gracefulShutdown(server, cancel) })
	eg.Go(server.ListenAndServe)

	kafkaReader := kafka.NewKafkaOrderConsumer(database.New())

	eg.Go(func() error { return kafkaReader.Consume(ctx) })

	slog.Info("Service started")

	if err := eg.Wait(); err != nil {
		slog.Error(err.Error())
	}
	slog.Info("Graceful shutdown complete.")
}

func setupLogger() {
	logLevel := new(slog.LevelVar)
	options := &slog.HandlerOptions{Level: logLevel}

	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel.Set(slog.LevelDebug)
		options.AddSource = true
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, options))
	slog.SetDefault(logger)
}
