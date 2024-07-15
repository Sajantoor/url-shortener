package store

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"go.uber.org/zap"

	cache "github.com/Sajantoor/url-shortener/services/common/store/cache"
	db "github.com/Sajantoor/url-shortener/services/common/store/database"
	types "github.com/Sajantoor/url-shortener/services/common/utils"
	"github.com/go-redis/redis"
	gocql "github.com/gocql/gocql"
)

type Store struct {
	redis     *cache.Redis
	cassandra *db.Cassandra
}

func New(ctx context.Context) *Store {
	zap.L().Info("Connecting to datastore...")

	return &Store{
		redis:     cache.New(ctx),
		cassandra: db.New(),
	}
}

func (s *Store) Close() {
	s.redis.Close()
	s.cassandra.Close()
}

func (s *Store) GetUrlMapping(ctx context.Context, shortUrl string) (*URLMapping, error) {
	cached, err := s.redis.GetClient(ctx).Get(shortUrl).Result()

	switch {
	case errors.Is(err, redis.Nil):
		// key doesn't exist, continue to fetch from database
	case err != nil:
		zap.L().Sugar().Infof("Error getting value from cache: ", err)
		return nil, types.InternalServerError(err)
	case cached == "":
		zap.L().Sugar().Info("Empty value found in cache")
	case true: // value found in cache
		cachedValue := &URLMapping{}
		err := json.Unmarshal([]byte(cached), cachedValue)

		if err != nil {
			zap.L().Sugar().Error("Failed to unmarshal cached value ", err)
			return nil, types.InternalServerError(err)
		}

		return cachedValue, nil
	}

	db := s.cassandra.Client()
	result := &URLMapping{}
	query := db.Query("SELECT * FROM url_shortener.url_map WHERE short_url = ?", shortUrl).WithContext(ctx)
	err = query.Scan(&result.ShortURL, &result.CreatedAt, &result.LongURL)

	switch {
	case errors.Is(err, gocql.ErrNotFound):
		return nil, types.NotFoundError(err)
	case err != nil:
		zap.L().Sugar().Error("Failed to get result from Cassandra: ", err)
		return nil, types.InternalServerError(err)
	}

	// set value in cache
	value, err := json.Marshal(result)
	if err != nil {
		zap.L().Sugar().Error("Failed to marshal value", err)
		return nil, types.InternalServerError(err)
	}

	s.redis.GetClient(ctx).Set(shortUrl, value, 0)
	return result, nil
}

func (s *Store) CreateUrlMapping(ctx context.Context, longUrl string, shortUrl string) (*URLMapping, error) {
	db := s.cassandra.Client()
	createdAt := time.Now()

	// Create long to short mapping first to avoid duplicates
	query := db.Query("INSERT INTO url_shortener.long_to_short (short_url, long_url) VALUES (?, ?) IF NOT EXISTS", shortUrl, longUrl).WithContext(ctx)
	applied, err := query.ScanCAS()

	if !applied {
		zap.L().Sugar().Info("Long URL already exists: ", longUrl)
		return nil, types.AlreadyExistsError(errors.New("long URL already exists"))
	}

	if err != nil {
		zap.L().Sugar().Error("Failed to insert into Cassandra long_to_short mapping: ", err)
		return nil, types.InternalServerError(err)
	}

	query = db.Query("INSERT INTO url_shortener.url_map (short_url, long_url, created_at) VALUES (?, ?, ?) IF NOT EXISTS", shortUrl, longUrl, createdAt).WithContext(ctx)
	err = query.Exec()

	if err != nil {
		zap.L().Sugar().Error("Failed to insert into Cassandra: ", err)
		return nil, types.InternalServerError(err)
	}

	return &URLMapping{
		ShortURL:  shortUrl,
		LongURL:   longUrl,
		CreatedAt: createdAt,
	}, nil
}

func (s *Store) DeleteURLMapping(shortURL string) error {
	return nil
}
