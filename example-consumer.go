package main

import (
	"context"
	"fmt"

	gomq_consumer "github.com/MatheusCoxxxta/go-bullmq-consumer/worker"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

func SendFirstMail(ctx context.Context, data map[string]any) error {
	fmt.Println("SendFirstMail to", data["email"])
	return nil
}

func CreateEmailWorker(redisClient *redis.Client) gomq_consumer.Worker {

	emailQueueWorker := gomq_consumer.Worker{
		Instance: redisClient,
		Queue:    "emailQueue",
		Handlers: gomq_consumer.Handlers{
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

func CreatePaymentWorker(redisClient *redis.Client) gomq_consumer.Worker {

	paymentQueueWorker := gomq_consumer.Worker{
		Instance: redisClient,
		Queue:    "paymentQueue",
		Handlers: gomq_consumer.Handlers{
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

	go gomq_consumer.StartWorker(emailQueueWorker)
	go gomq_consumer.StartWorker(paymentQueueWorker)

	select {}
}
