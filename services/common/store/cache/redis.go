package cache

import (
	"log"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func New() *Redis {
	log.Println("Creating new Redis client")

	client := redis.NewClient(&redis.Options{
		// TODO: Add additional configuration here
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()

	if err != nil {
		log.Fatal("Failed to connect to Redis", err)
	}

	log.Println("Connected to Redis")
	return &Redis{
		client: client,
	}
}

func (r *Redis) Close() {
	r.client.Close()
}

func (r *Redis) GetClient() *redis.Client {
	return r.client
}
