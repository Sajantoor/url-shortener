package handler

import (
	"context"
	"net/url"
	"os"

	pb "github.com/Sajantoor/url-shortener/services/common/protobuf"
	"github.com/Sajantoor/url-shortener/services/common/store"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RetervialHandler struct {
	Store *store.Store
	pb.UnimplementedUrlRetrievalServer
}

var domain = os.Getenv("DOMAIN")
var domainHostname string = ""

func (s *RetervialHandler) GetLongUrl(ctx context.Context, req *pb.GetLongUrlRequest) (*pb.GetLongUrlResponse, error) {
	if req.GetShortUrl() == "" {
		return nil, status.Error(codes.InvalidArgument, "Short URL cannot be empty")
	}

	shortUrl, err := url.Parse(req.GetShortUrl())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid URL, please provide a valid URL")
	}

	shortUrlHostName := shortUrl.Hostname()
	if shortUrlHostName == "" {
		return nil, status.Error(codes.InvalidArgument, "Invalid short URL")
	}

	domainHostName := getDomainHostName()

	if shortUrlHostName != domainHostName {
		return nil, status.Error(codes.InvalidArgument, "Invalid short URL")
	}

	shortUrlHash := shortUrl.Path

	if shortUrlHash == "" && len(shortUrlHash) <= 1 {
		return nil, status.Error(codes.InvalidArgument, "Invalid short URL")
	}

	shortUrlHash = shortUrlHash[1:]

	shortUrlMapping, err := s.Store.GetUrlMapping(shortUrlHash)

	if err != nil {
		return nil, status.Error(codes.NotFound, "Short URL not found")
	}

	return &pb.GetLongUrlResponse{
		LongUrl:   shortUrlMapping.LongURL,
		CreatedAt: timestamppb.New(shortUrlMapping.CreatedAt),
	}, nil
}

func getDomainHostName() string {
	if domainHostname != "" {
		return domainHostname
	}

	parsedDomain, err := url.Parse(domain)
	if err != nil {
		panic("Invalid domain")
	}

	domainHostname = parsedDomain.Hostname()
	return domainHostname
}
