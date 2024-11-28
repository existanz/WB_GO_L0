package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"WB_GO_L0/internal/kafka"
	"WB_GO_L0/internal/server"

	"golang.org/x/sync/errgroup"
)

func gracefulShutdown(apiServer *http.Server, cancelKafka context.CancelFunc) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	cancelKafka()
	log.Println("Server exiting")

	return nil
}

func main() {
	server := server.NewServer()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eg := new(errgroup.Group)

	eg.Go(func() error { return gracefulShutdown(server, cancel) })
	eg.Go(server.ListenAndServe)
	eg.Go(func() error { return kafka.Consume(ctx) })

	if err := eg.Wait(); err != nil {
		fmt.Println(err)
	}
	log.Println("Graceful shutdown complete.")
}
