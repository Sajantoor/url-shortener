package main

import (
	"context"
	"os"

	"github.com/Sajantoor/url-shortener/services/common/grpcServer"
	"github.com/Sajantoor/url-shortener/services/common/protobuf"
	"github.com/Sajantoor/url-shortener/services/common/store"
	"github.com/Sajantoor/url-shortener/services/common/utils"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/Sajantoor/url-shortener/services/retrieval/handler"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	utils.InitLogger()
	zap.L().Info("Starting URL Shortener Retervial Service...")
	godotenv.Load("../common/.env")

	port := os.Getenv("RETRIEVAL_SERVICE_PORT")
	if port == "" {
		panic("RETRIEVAL_SERVICE_PORT is not set")
	}

	store := store.New(ctx)
	defer store.Close()

	grpcServer := grpcServer.New(":" + port)

	healthServer := utils.NewHealthServer()
	healthpb.RegisterHealthServer(grpcServer.GetServer(), healthServer)

	protobuf.RegisterUrlRetrievalServer(grpcServer.GetServer(), &handler.RetervialHandler{
		Store: store,
	})

	healthServer.SetStatus("", healthpb.HealthCheckResponse_SERVING)

	go func() {
		grpcServer.Start()
	}()

	utils.HandleShutdown(ctx, healthServer)
}
