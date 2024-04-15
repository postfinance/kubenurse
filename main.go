// Package main holds the entrypoint for kubenurse
package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/postfinance/kubenurse/internal/kubenurse"
	corev1 "k8s.io/api/core/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	restConf, err := controllerruntime.GetConfig()
	if err != nil {
		slog.Error("error during controllerruntime.GetConfig()", "err", err)
		return
	}

	kubenurseNs := os.Getenv("KUBENURSE_NAMESPACE")

	ca, err := cache.New(restConf, cache.Options{
		ByObject: map[client.Object]cache.ByObject{
			&corev1.Pod{}: {Namespaces: map[string]cache.Config{
				kubenurseNs: {},
			}},
			&corev1.Node{}: {},
		},
	})

	if err != nil {
		slog.Error("error during cache creation", "err", err)
		return
	}

	go func() {
		if err = ca.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			slog.Error("controller-runtime client cache error", "err", err)
			cancel()
		}
	}()

	opts := client.Options{
		Cache: &client.CacheOptions{
			Reader: ca,
		},
	}

	c, err := client.New(restConf, opts)
	if err != nil {
		slog.Error("error while starting controller-runtime client", "err", err)
		return
	}

	server, err := kubenurse.New(ctx, c)
	if err != nil {
		slog.Error("error in kubenurse.New call", "err", err)
		return
	}

	go func() {
		<-ctx.Done() // blocks until ctx is canceled

		slog.Info("shutting down, received signal to stop")

		// background ctx since, the "root" context is already canceled
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("error during graceful shutdown", "err", err)
		}
	}()

	// blocks, until the server is stopped by calling Shutdown()
	if err := server.Run(); err != nil {
		slog.Error("error while running kubenurse", "err", err)
	}
}
