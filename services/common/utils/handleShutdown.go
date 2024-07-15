package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func HandleShutdown(ctx context.Context, healthServer *HealthServer) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		zap.L().Info("Shutting down " + v.String())
		healthServer.SetStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	case <-ctx.Done():
		zap.L().Info("Shutting down")
		healthServer.SetStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	}
}
