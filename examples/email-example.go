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

	emailQueueWorker := oncamq.New(
		ctx,
		redisClient,
		"emailQueue",
		oncamq.Handlers{
			"firstAccess": SendFirstMail,
		},
	)

	fmt.Println("Starting Email Queue Worker...")
	emailQueueWorker.Start()

	// Keep the main function running
	emailQueueWorker.Wait()
}
