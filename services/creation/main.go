package main

import (
	"fmt"

	"github.com/Sajantoor/url-shortener/services/common/grpcServer"
	"github.com/Sajantoor/url-shortener/services/common/protobuf"
	Cassandra "github.com/Sajantoor/url-shortener/services/common/store"
	"github.com/Sajantoor/url-shortener/services/creation/handler"
)

func main() {
	fmt.Println("Starting URL Shortener Creation Service...")

	fmt.Println("Connecting to database...")

	Cassandra.New()

	grpcServer := grpcServer.New(":8080")

	protobuf.RegisterUrlShortnerServiceServer(grpcServer.GetServer(), &handler.CreationHandler{})

	grpcServer.Start()
}
