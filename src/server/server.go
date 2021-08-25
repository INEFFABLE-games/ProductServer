package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"log"
	"main/src/internal/models"
	"time"
)

func getRedis() *redis.Client {
	var (
		RedisAddres    = "redis-14436.c54.ap-northeast-1-2.ec2.cloud.redislabs.com:14436"
		RedisPassword  = "Drektarov3698!"
		ReddisUserName = "Admin"
	)

	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddres,
		Password: RedisPassword,
		Username: ReddisUserName,
		DB:       0,
	})

	res, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)

	return client
}

func Start(ctx context.Context) {
	topic := "products"

	//ip := os.Getenv("ZADDRES")

	conf := kafka.ReaderConfig{
		Brokers:  []string{"localhost:29092"},
		Topic:    topic,
		GroupID:  "g1",
		MaxBytes: 10,
	}

	redisClient := getRedis()

	reader := kafka.NewReader(conf)

	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			fmt.Println("Some error: " + err.Error())
			continue
		}

		prod := models.Product{}

		err = json.Unmarshal(m.Value, &prod)
		if err != nil {
			fmt.Println(err)
		}
		jprod, err := json.Marshal(prod)
		if err != nil {
			log.Println(err)
		}
		redisClient.Set(ctx, "product", jprod, 1*time.Second)
		redisClient.Set(ctx, "product1", jprod, 1*time.Second)

		fmt.Printf("new product: %v\n", prod)

		time.Sleep(3 * time.Second)
	}

}
