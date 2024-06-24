package main

import (
	"log"
	"net"

	"github.com/Sajantoor/url-shortener/services/common/protobuf"
	"github.com/Sajantoor/url-shortener/services/creation/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to create TCP listener:", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	protobuf.RegisterUrlShortnerServiceServer(grpcServer, &server.CreationServer{})
	log.Print("Starting gRPC server on port 8080")

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalln("Failed to start gRPC server:", err)
	}

	log.Print("gRPC server started successfully")
}
