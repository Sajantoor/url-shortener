package server

import (
	"context"
	"time"

	"github.com/Sajantoor/url-shortener/services/common/protobuf"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreationServer struct {
	protobuf.UnimplementedUrlShortnerServiceServer
}

func (s *CreationServer) CreateShortUrl(ctx context.Context, req *protobuf.CreateShortUrlRequest) (*protobuf.CreateShortUrlResponse, error) {
	return &protobuf.CreateShortUrlResponse{
		ShortUrl:  "http://localhost:3000/short/",
		CreatedAt: timestamppb.New(time.Now()),
	}, nil
}
