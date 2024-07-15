package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func HandleShutdown(ctx context.Context) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		zap.L().Info("Shutting down " + v.String())
	case <-ctx.Done():
		zap.L().Info("Shutting down")
	}
}
