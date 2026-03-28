package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx    = context.Background()
	Client *redis.Client
)

func Init() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	for i := 0; i < 5; i++ {
		_, err := client.Ping(Ctx).Result()

		if err == nil {
			log.Println("Connected to Redis")
			return
		}
		log.Println("Waiting for Redis...")
		time.Sleep(2 * time.Second)
	}

	log.Fatal("Could not connect to Redis")
}
