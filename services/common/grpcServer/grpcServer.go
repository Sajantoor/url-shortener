package grpcServer

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	listener   *net.Listener
	grpcServer *grpc.Server
}

func New(address string) *GrpcServer {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatalln("Failed to create TCP listener:", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	return &GrpcServer{
		listener:   &listener,
		grpcServer: grpcServer,
	}
}

func (s *GrpcServer) Start() {
	err := s.grpcServer.Serve(*s.listener)

	if err != nil {
		log.Fatalln("Failed to start gRPC server:", err)
	}

	log.Print("gRPC server started successfully")
}

func (s *GrpcServer) GetServer() *grpc.Server {
	return s.grpcServer
}
