package main

import (
	"log"
	"os"

	"github.com/Sajantoor/url-shortener/services/common/grpcServer"
	"github.com/Sajantoor/url-shortener/services/common/protobuf"
	"github.com/Sajantoor/url-shortener/services/common/store"
	"github.com/Sajantoor/url-shortener/services/creation/handler"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../common/.env")
	log.Println("Starting URL Shortener Creation Service...")

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
