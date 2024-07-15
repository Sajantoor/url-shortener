package main

import (
	"os"

	"github.com/Sajantoor/url-shortener/services/common/grpcServer"
	"github.com/Sajantoor/url-shortener/services/common/protobuf"
	"github.com/Sajantoor/url-shortener/services/common/store"
	logger "github.com/Sajantoor/url-shortener/services/common/types"

	"github.com/Sajantoor/url-shortener/services/retrieval/handler"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	godotenv.Load("../common/.env")

	zap.L().Info("Starting URL Shortener Retervial Service...")

	store := store.New()
	defer store.Close()

	port := os.Getenv("RETRIEVAL_SERVICE_PORT")
	if port == "" {
		panic("RETRIEVAL_SERVICE_PORT is not set")
	}

	grpcServer := grpcServer.New(":" + port)

	protobuf.RegisterUrlRetrievalServer(grpcServer.GetServer(), &handler.RetervialHandler{
		Store: store,
	})

	grpcServer.Start()
}
