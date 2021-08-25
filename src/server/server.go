package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx"
	"github.com/segmentio/kafka-go"
	"log"
	"main/src/internal/models"
	"main/src/internal/repository"
	"main/src/internal/services"
	"os"
	"time"
	//"github.com/jackc/pgx"
)

func getPostgre() *pgx.Conn {
	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:                 "localhost",
		Port:                 5432,
		Database:             "products",
		User:                 "postgres",
		Password:             "12345",
		TLSConfig:            nil,
		UseFallbackTLS:       false,
		FallbackTLSConfig:    nil,
		Logger:               nil,
		LogLevel:             0,
		Dial:                 nil,
		RuntimeParams:        nil,
		OnNotice:             nil,
		CustomConnInfo:       nil,
		CustomCancel:         nil,
		PreferSimpleProtocol: false,
		TargetSessionAttrs:   "",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}

	return conn
}

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
	//-------------------------------------------------------POSTGRESQL

	conn := getPostgre()
	productPostgreRepo := repository.NewPostgreProductRepository(conn)
	productService := services.NewProductService(&productPostgreRepo)
	defer conn.Close()

	//---------------------------------------------------------KAFKA
	conf := kafka.ReaderConfig{
		Brokers:  []string{"localhost:29092"},
		Topic:    topic,
		GroupID:  "g1",
		MaxBytes: 10,
	}

	redisClient := getRedis() //-----------------------------redis

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

		err = redisClient.Do(ctx, "XADD", "products", "*", "name", prod.Name, "price", prod.Price).Err()
		if err != nil {
			log.Println(err)
		}

		err = productService.AddProduct(prod)
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("new product: %v\n", prod)

		time.Sleep(3 * time.Second)
	}

}
