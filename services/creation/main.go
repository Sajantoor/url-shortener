package main

import (
	"log"

	"github.com/Sajantoor/url-shortener/services/common/grpcServer"
	"github.com/Sajantoor/url-shortener/services/common/protobuf"
	"github.com/Sajantoor/url-shortener/services/common/store"
	"github.com/Sajantoor/url-shortener/services/creation/handler"
)

func main() {
	log.Println("Starting URL Shortener Creation Service...")

	store := store.New()
	defer store.Close()

	grpcServer := grpcServer.New(":8080")

	protobuf.RegisterUrlShortnerServiceServer(grpcServer.GetServer(), &handler.CreationHandler{
		Store: store,
	})

	grpcServer.Start()
}
