package handler

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/url"
	"os"

	"github.com/Sajantoor/url-shortener/services/common/protobuf"
	"github.com/Sajantoor/url-shortener/services/common/store"
	types "github.com/Sajantoor/url-shortener/services/common/utils"
	"google.golang.org/protobuf/types/known/timestamppb"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var domain = os.Getenv("DOMAIN")

const shortUrlLength = 7
const hashLength = 5
const randomStrLength = shortUrlLength - hashLength

type CreationHandler struct {
	Store *store.Store
	protobuf.UnimplementedUrlShortnerServiceServer
}

func (s *CreationHandler) CreateShortUrl(ctx context.Context, req *protobuf.CreateShortUrlRequest) (*protobuf.CreateShortUrlResponse, error) {
	longUrl := req.GetLongUrl()

	if longUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "Long URL cannot be empty")
	}

	// validate it is a valid url
	_, err := url.ParseRequestURI(longUrl)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid URL, please provide a valid URL")
	}

	shortUrl := generateShortUrl(longUrl)

	res, err := s.Store.CreateUrlMapping(ctx, longUrl, shortUrl)

	if err != nil {
		if errors.As(err, &types.ReqError) {
			// cast to ReqError and return the error message
			reqErr := err.(*types.RequestError)
			return nil, status.Error(reqErr.Code, reqErr.Error())
		}

		return nil, status.Error(codes.Internal, "Failed to create short URL")
	}

	return &protobuf.CreateShortUrlResponse{
		ShortUrl:  domain + res.ShortURL,
		CreatedAt: timestamppb.New(res.CreatedAt),
	}, nil
}

func generateShortUrl(longUrl string) string {
	hash := generateHash(longUrl)
	randomStr := generateRandomString(shortUrlLength)
	shortUrl := hash[:hashLength] + randomStr[:randomStrLength]

	// TODO: Check if this exists in the database...

	return shortUrl
}

func generateHash(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length]
}
