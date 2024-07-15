package grpcServer

import (
	"net"

	"go.uber.org/zap"

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
		zap.L().Fatal("Failed to create TCP listener:" + err.Error())
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
		zap.L().Fatal("Failed to start gRPC server:" + err.Error())
	}

	zap.L().Sugar().Info("gRPC server started successfully")
}

func (s *GrpcServer) GetServer() *grpc.Server {
	return s.grpcServer
}
