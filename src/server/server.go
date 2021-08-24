package server

import (
	"context"
)

func Start(ctx context.Context) {
	/*topic := "my-topic"

	//ip := os.Getenv("ZADDRES")

	conf := kafka.ReaderConfig{
		Brokers:  []string{"localhost:29092"},
		Topic:    topic,
		GroupID:  "g1",
		MaxBytes: 10,
	}

	reader := kafka.NewReader(conf)

	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			fmt.Println("Some error: " + err.Error())
			continue
		}

		err = json.Unmarshal(m.Value,&models.Product{}); if err != nil{
			fmt.Println(err)
		}

		time.Sleep(1 * time.Second)
	}*/

}
