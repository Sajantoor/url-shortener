package store

import (
	"encoding/json"
	"log"
	"time"

	cache "github.com/Sajantoor/url-shortener/services/common/store/cache"
	db "github.com/Sajantoor/url-shortener/services/common/store/database"
)

type Store struct {
	redis     *cache.Redis
	cassandra *db.Cassandra
}

type UrlAlreadyExistsError struct{}

func (z UrlAlreadyExistsError) Error() string {
	return "URL already exists"
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

	if err != nil {
		// return value from cache
		cachedValue := &URLMapping{}

		err := json.Unmarshal([]byte(cached), cachedValue)

		if err != nil {
			log.Fatal("Failed to unmarshal cached value", err)
			// TODO: Maybe we swallow this error
			return nil, err
		}

		return cachedValue, nil
	}

	db := s.cassandra.Client()
	result := &URLMapping{}
	query := db.Query("SELECT * FROM url_shortener.url_map WHERE short_url = ?", shortUrl)

	err = query.Scan(&result.ShortURL, &result.CreatedAt, &result.LongURL)

	if err != nil {
		log.Println("Error getting result: ", err)
	}

	// set value in cache
	value, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Failed to marshal value", err)
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
		// TODO: Improve this error...
		return nil, UrlAlreadyExistsError{}
	}

	if err != nil {
		log.Println("Failed to insert into Cassandra: ", err)
		return nil, err
	}

	query = db.Query("INSERT INTO url_shortener.url_map (short_url, long_url, created_at) VALUES (?, ?, ?) IF NOT EXISTS", shortUrl, longUrl, createdAt)

	err = query.Exec()

	if err != nil {
		log.Println("Failed to insert into Cassandra: ", err)
		return nil, err
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
