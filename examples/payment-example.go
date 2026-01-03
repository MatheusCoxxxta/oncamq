package main

import (
	"context"
	"fmt"

	"github.com/MatheusCoxxxta/oncamq"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

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
	ctx := context.Background()
	redisClient := InitRedis()

	paymentQueueWorker := oncamq.New(
		ctx,
		redisClient,
		"paymentQueue",
		oncamq.Handlers{
			"createCustomer":   CreateCustomer,
			"startTransaction": StartTransaction,
		},
	)

	fmt.Println("Starting Payment Queue Worker...")
	paymentQueueWorker.Start()
	
	// Keep the main function running
	paymentQueueWorker.Wait()
}
