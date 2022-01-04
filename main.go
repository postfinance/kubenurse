package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/postfinance/kubenurse/internal/kubenurse"
)

const (
	nurse = "I'm ready to help you!"
)

func main() {
	server, err := kubenurse.New()
	if err != nil {
		log.Printf("%s", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done() // blocks until ctx is canceled

		log.Println("shutting down, received signal to stop")

		// background ctx since, the "root" context is already canceled
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("gracefully halting kubenurse server: %s", err)
		}
	}()

	fmt.Println(nurse) // most important line of this project

	// blocks, until the server is stopped by calling Shutdown()
	if err := server.Run(); err != nil {
		log.Printf("running kubenurse: %s", err)
	}
}
