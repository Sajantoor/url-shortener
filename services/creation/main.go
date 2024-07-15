package main

import (
	"os"

	"github.com/Sajantoor/url-shortener/services/common/grpcServer"
	"github.com/Sajantoor/url-shortener/services/common/protobuf"
	"github.com/Sajantoor/url-shortener/services/common/store"
	logger "github.com/Sajantoor/url-shortener/services/common/types"
	"github.com/Sajantoor/url-shortener/services/creation/handler"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	godotenv.Load("../common/.env")

	zap.L().Info("Starting URL Shortener Creation Service...")

	store := store.New()
	defer store.Close()

	port := os.Getenv("CREATION_SERVICE_PORT")
	if port == "" {
		panic("CREATION_SERVICE_PORT is not set")
	}

	grpcServer := grpcServer.New(":" + port)

	protobuf.RegisterUrlShortnerServiceServer(grpcServer.GetServer(), &handler.CreationHandler{
		Store: store,
	})

	grpcServer.Start()
}
