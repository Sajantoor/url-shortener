package main

import (
	"log"
	"os"

	"github.com/Sajantoor/url-shortener/services/common/grpcServer"
	"github.com/Sajantoor/url-shortener/services/common/protobuf"
	"github.com/Sajantoor/url-shortener/services/common/store"
	"github.com/Sajantoor/url-shortener/services/retrieval/handler"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../common/.env")
	log.Println("Starting URL Shortener Retervial Service...")

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
