// Package main holds the entrypoint for kubenurse
package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/postfinance/kubenurse/internal/kubenurse"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	nurse = "I'm ready to help you!"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// create in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("creating in-cluster configuration: %s", err)
		return
	}

	cliset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("creating clientset: %s", err)
		return
	}

	server, err := kubenurse.New(ctx, cliset)
	if err != nil {
		log.Printf("%s", err)
		return
	}

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
