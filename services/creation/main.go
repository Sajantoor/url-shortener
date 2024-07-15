package main

import (
	"context"
	"os"

	"github.com/Sajantoor/url-shortener/services/common/grpcServer"
	"github.com/Sajantoor/url-shortener/services/common/protobuf"
	"github.com/Sajantoor/url-shortener/services/common/store"
	"github.com/Sajantoor/url-shortener/services/common/utils"
	"github.com/Sajantoor/url-shortener/services/creation/handler"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	utils.InitLogger()
	godotenv.Load()

	zap.L().Info("Starting URL Shortener Creation Service...")

	port := os.Getenv("CREATION_SERVICE_PORT")
	if port == "" {
		panic("CREATION_SERVICE_PORT is not set")
	}
	zap.L().Info("Starting application on port: " + port)

	store := store.New(ctx)
	defer store.Close()

	grpcServer := grpcServer.New(":" + port)

	healthServer := utils.NewHealthServer()
	healthpb.RegisterHealthServer(grpcServer.GetServer(), healthServer)

	protobuf.RegisterUrlShortnerServiceServer(grpcServer.GetServer(), &handler.CreationHandler{
		Store: store,
	})

	healthServer.SetStatus("", healthpb.HealthCheckResponse_SERVING)

	go func() {
		grpcServer.Start()
	}()

	utils.HandleShutdown(ctx, healthServer)
}
