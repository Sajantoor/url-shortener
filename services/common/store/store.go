package store

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	cache "github.com/Sajantoor/url-shortener/services/common/store/cache"
	db "github.com/Sajantoor/url-shortener/services/common/store/database"
	"github.com/Sajantoor/url-shortener/services/common/types"
	"github.com/go-redis/redis"
	gocql "github.com/gocql/gocql"
)

type Store struct {
	redis     *cache.Redis
	cassandra *db.Cassandra
}

func New() *Store {
	log.Println("Connecting to datastore...")

	return &Store{
		redis:     cache.New(),
		cassandra: db.New(),
	}
}

func (s *Store) Close() {
	s.redis.Close()
	s.cassandra.Close()
}

func (s *Store) GetUrlMapping(shortUrl string) (*URLMapping, error) {
	cached, err := s.redis.GetClient().Get(shortUrl).Result()

	switch {
	case errors.Is(err, redis.Nil):
		// key doesn't exist, continue to fetch from database
	case err != nil:
		log.Println("Error getting value from cache: ", err)
		return nil, types.InternalServerError(err)
	case cached == "":
		log.Println("Empty value found in cache")
	case true: // value found in cache
		cachedValue := &URLMapping{}
		err := json.Unmarshal([]byte(cached), cachedValue)

		if err != nil {
			log.Fatal("Failed to unmarshal cached value ", err)
			return nil, types.InternalServerError(err)
		}

		return cachedValue, nil
	}

	db := s.cassandra.Client()
	result := &URLMapping{}
	query := db.Query("SELECT * FROM url_shortener.url_map WHERE short_url = ?", shortUrl)
	err = query.Scan(&result.ShortURL, &result.CreatedAt, &result.LongURL)

	switch {
	case errors.Is(err, gocql.ErrNotFound):
		return nil, types.NotFoundError(err)
	case err != nil:
		log.Println("Failed to get result from Cassandra: ", err)
		return nil, types.InternalServerError(err)
	}

	// set value in cache
	value, err := json.Marshal(result)
	if err != nil {
		log.Println("Failed to marshal value", err)
		return nil, types.InternalServerError(err)
	}

	s.redis.GetClient().Set(shortUrl, value, 0)
	return result, nil
}

func (s *Store) CreateUrlMapping(longUrl string, shortUrl string) (*URLMapping, error) {
	db := s.cassandra.Client()
	createdAt := time.Now()

	// Create long to short mapping first to avoid duplicates
	query := db.Query("INSERT INTO url_shortener.long_to_short (short_url, long_url) VALUES (?, ?) IF NOT EXISTS", shortUrl, longUrl)
	applied, err := query.ScanCAS()

	if !applied {
		log.Println("Long URL already exists: ", longUrl)
		return nil, types.AlreadyExistsError(errors.New("long URL already exists"))
	}

	if err != nil {
		log.Println("Failed to insert into Cassandra long_to_short mapping: ", err)
		return nil, types.InternalServerError(err)
	}

	query = db.Query("INSERT INTO url_shortener.url_map (short_url, long_url, created_at) VALUES (?, ?, ?) IF NOT EXISTS", shortUrl, longUrl, createdAt)
	err = query.Exec()

	if err != nil {
		log.Println("Failed to insert into Cassandra: ", err)
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
