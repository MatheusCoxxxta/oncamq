package main

import (
	"context"
	"fmt"

	"github.com/MatheusCoxxxta/go-bullmq-consumer/worker"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

func SendFirstMail(ctx context.Context, data map[string]any) error {
	fmt.Println("SendFirstMail to", data["email"])
	return nil
}

func CreateEmailWorker(redisClient *redis.Client) worker.Worker {

	emailQueueWorker := worker.Worker{
		Instance: redisClient,
		Queue:    "emailQueue",
		Handlers: worker.Handlers{
			"firstAcess": SendFirstMail,
		},
	}

	return emailQueueWorker
}

func CreateCustomer(ctx context.Context, data map[string]any) error {
	fmt.Println("CreateCustomer to", data["email"])
	return nil
}

func StartTransaction(ctx context.Context, data map[string]any) error {
	fmt.Println("StartTransaction to", data["email"])
	return nil
}

func CreatePaymentWorker(redisClient *redis.Client) worker.Worker {

	paymentQueueWorker := worker.Worker{
		Instance: redisClient,
		Queue:    "paymentQueue",
		Handlers: worker.Handlers{
			"createCustomer":   CreateCustomer,
			"startTransaction": StartTransaction,
		},
	}

	return paymentQueueWorker
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

	emailQueueWorker := CreateEmailWorker(redisClient)
	paymentQueueWorker := CreatePaymentWorker(redisClient)

	go worker.StartWorker(emailQueueWorker)
	go worker.StartWorker(paymentQueueWorker)

	select {}
}
