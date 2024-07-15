package main

import (
	"context"
	"os"

	"github.com/Sajantoor/url-shortener/services/common/grpcServer"
	"github.com/Sajantoor/url-shortener/services/common/protobuf"
	"github.com/Sajantoor/url-shortener/services/common/store"
	"github.com/Sajantoor/url-shortener/services/common/utils"

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

	protobuf.RegisterUrlRetrievalServer(grpcServer.GetServer(), &handler.RetervialHandler{
		Store: store,
	})

	go func() {
		grpcServer.Start()
	}()

	utils.HandleShutdown(ctx)
}
