package main

import (
	"context"
	"fmt"

	"github.com/MatheusCoxxxta/oncamq"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

func SendFirstMail(ctx context.Context, data map[string]any) (any, error) {
	fmt.Println("SendFirstMail to", data["to"])

	return "Sent", nil
}

func CreateCustomer(ctx context.Context, data map[string]any) (any, error) {
	fmt.Println("CreateCustomer to", data["email"])

	return "Created", nil
}

func StartTransaction(ctx context.Context, data map[string]any) (any, error) {
	fmt.Println("StartTransaction to", data["email"])

	return "Started", nil
}

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6370",
		Password: "redis",
		DB:       0,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})

	return client
}

func main() {
	redisClient := InitRedis()
	ctx := context.Background()

	emailQueueWorker := oncamq.New(
		redisClient,
		"emailQueue",
		oncamq.Handlers{
			"firstAccess": SendFirstMail,
		},
	)
	paymentQueueWorker := oncamq.New(
		redisClient,
		"paymentQueue",
		oncamq.Handlers{
			"createCustomer":   CreateCustomer,
			"startTransaction": StartTransaction,
		},
	)

	go emailQueueWorker.StartWorker(ctx)
	go paymentQueueWorker.StartWorker(ctx)

	select {}
}
