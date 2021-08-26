package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"github.com/segmentio/kafka-go"
	"log"
	"main/src/internal/models"
	"main/src/internal/repository"
	"main/src/internal/services"
	"time"
	//"github.com/jackc/pgx"
)

func getPostgre() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), "port=5432 host=localhost user=postgres password=12345 dbname=products sslmode=disable")
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
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
	defer conn.Close(ctx)

	//---------------------------------------------------------KAFKA
	conf := kafka.ReaderConfig{
		Brokers:        []string{"localhost:29092"},
		Topic:          topic,
		GroupID:        "g1",
		MaxBytes:       10,
		CommitInterval: 10 * time.Second,
	}

	//-----------------------------redis
	redisClient := getRedis()
	err := redisClient.Do(ctx, "DEL", "products").Err()
	if err != nil {
		log.Println(err)
	}

	//----------------------------kafka
	reader := kafka.NewReader(conf)

	for {

		var newproducts = make([]models.Product, 10)

		for i := 0; i < 10; i++ {
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
			newproducts[i] = prod

			err = redisClient.Do(ctx, "XADD", "products", "*", "name", prod.Name, "price", prod.Price).Err()
			if err != nil {
				log.Println(err)
			}

			fmt.Printf("new product: %v\n", prod)
		}

		err = productService.AddProduct(newproducts)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(3 * time.Second)
	}

}
