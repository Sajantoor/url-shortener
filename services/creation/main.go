package main

import (
	"github.com/Sajantoor/url-shortener/services/common/grpcServer"
	"github.com/Sajantoor/url-shortener/services/common/protobuf"
	"github.com/Sajantoor/url-shortener/services/creation/handler"
)

func main() {
	grpcServer := grpcServer.New(":8080")

	protobuf.RegisterUrlShortnerServiceServer(grpcServer.GetServer(), &handler.CreationHandler{})

	grpcServer.Start()
}
