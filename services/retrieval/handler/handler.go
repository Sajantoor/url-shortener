package handler

import (
	"context"
	"errors"
	"net/url"

	pb "github.com/Sajantoor/url-shortener/services/common/protobuf"
	"github.com/Sajantoor/url-shortener/services/common/store"
	types "github.com/Sajantoor/url-shortener/services/common/utils"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RetervialHandler struct {
	Store *store.Store
	pb.UnimplementedUrlRetrievalServer
}

var domain = "http://localhost:8080"
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
		return nil, status.Error(codes.InvalidArgument, "Invalid short URL: Short URL cannot be empty.")
	}

	domainHostName := getDomainHostName()

	if shortUrlHostName != domainHostName {
		return nil, status.Error(codes.InvalidArgument, "Invalid short URL: Short URL does not belong to this service.")
	}

	shortUrlHash := shortUrl.Path

	if shortUrlHash == "" && len(shortUrlHash) <= 1 {
		return nil, status.Error(codes.InvalidArgument, "Invalid short URL: Short URL hash cannot be empty.")
	}

	shortUrlHash = shortUrlHash[1:]

	shortUrlMapping, err := s.Store.GetUrlMapping(ctx, shortUrlHash)

	if err != nil {
		if errors.As(err, &types.ReqError) {
			// cast to ReqError and return the error message
			reqErr := err.(*types.RequestError)
			return nil, status.Error(reqErr.Code, reqErr.Error())
		}

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
