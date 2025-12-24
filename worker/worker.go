package gomq_consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var redisClient *redis.Client

type Handlers map[string]func(ctx context.Context, data map[string]any) (any, error)

type Worker struct {
	Instance *redis.Client
	Queue    string
	Handlers Handlers
}

type GetKeyFromQueueResponse struct {
	Name        *string
	Data        *string
	IsProcessed bool
	IsNotFound  bool
}

func registerInstance(connection *redis.Client) {
	redisClient = connection
}

func createKey(queue string, key string) string {
	return fmt.Sprintf("bull:%s:%s", queue, key)
}

func getValueByKey(key string) GetKeyFromQueueResponse {

	job, err := redisClient.HGetAll(ctx, key).Result()

	if err != nil {
		fmt.Println(err)
	}

	if job["name"] == "" {
		return GetKeyFromQueueResponse{
			IsNotFound:  true,
			IsProcessed: false,
		}
	}

	if job["processedOn"] != "" {
		return GetKeyFromQueueResponse{
			IsProcessed: true,
			IsNotFound:  false,
		}
	}

	value, err := redisClient.HGet(ctx, key, "data").Result()

	if err != nil {
		fmt.Println(err)
	}

	jobName := job["name"]

	return GetKeyFromQueueResponse{
		Data:        &value,
		Name:        &jobName,
		IsProcessed: false,
		IsNotFound:  false,
	}
}

func createLastReadKey(queue string) string {
	return fmt.Sprintf("bull:%s:completed", queue)
}

func setJobProcessed(queue string, key int) {
	fullKeyJob := fmt.Sprintf("bull:%s:%d", queue, key)

	redisClient.HSet(ctx, fullKeyJob,
		"processedOn", time.Now().UnixMilli(),
	)
}

func setJobFinished(queue string, key int, finishedOn int64, returnValue interface{}) {
	fullKeyJob := fmt.Sprintf("bull:%s:%d", queue, key)

	redisClient.HSet(ctx, fullKeyJob,
		"finishedOn", finishedOn,
		"returnvalue", returnValueToString(returnValue),
	)
}

func setAttemptsCount(queue string, key int, attempt int, attemptType string) {

	fullKeyJob := fmt.Sprintf("bull:%s:%d", queue, key)

	redisClient.HSet(ctx, fullKeyJob,
		attemptType, attempt,
	)
}

func getLastIndexRead(queue string) int {

	fullKey := createLastReadKey(queue)

	value, err := redisClient.ZRange(ctx, fullKey, 0, -1).Result()

	if err != nil {
		fmt.Println(err)
	}

	if len(value) == 0 {
		return 0
	}

	integerValue, err := strconv.Atoi(value[len(value)-1])

	if err != nil {
		fmt.Println(err)
	}

	return integerValue
}

func getAttemptsCount(queue string, key int, attemptType string) int {

	fullKeyJob := fmt.Sprintf("bull:%s:%d", queue, key)

	value, err := redisClient.HGet(ctx, fullKeyJob, attemptType).Result()

	if err != nil {
		fmt.Println(err)
	}

	if len(value) == 0 {
		return 0
	}

	integerValue, err := strconv.Atoi(value)

	if err != nil {
		fmt.Println(err)
	}

	return integerValue
}

func addToJobCompletedQueue(queue string, key int, finishedOn int64) {
	fullKeyJob := fmt.Sprintf("bull:%s:completed", queue)

	redisClient.ZAdd(ctx, fullKeyJob, redis.Z{
		Score:  float64(finishedOn),
		Member: key,
	})
}

// To finish/solve previous event nuance
// func addEventForJobCompleted(queue string, key int, returnValue string, previousEvent string) {
// 	fullKeyJob := fmt.Sprintf("bull:%s:events", queue)

// 	redisClient.XAdd(ctx, &redis.XAddArgs{
// 		Stream: fullKeyJob,
// 		ID:     "*",
// 		Values: map[string]interface{}{
// 			"event":       "completed",
// 			"jobId":       strconv.Itoa(key),
// 			"returnValue": returnValue,
// 			"prev":        previousEvent,
// 		},
// 	})
// }

func returnValueToString(value any) string {
	if value == nil {
		return "null"
	}

	jsonEncoding, err := json.Marshal(value)

	if err != nil {
		return ""
	}

	return string(jsonEncoding)
}

func StartWorker(worker Worker) {
	log.Printf("Worker to consume %s (%s) started...", worker.Queue, worker.Instance.Options().Addr)

	registerInstance(worker.Instance)

	for {
		lastRead := getLastIndexRead(worker.Queue)

		currentIndexToRead := lastRead + 1
		key := createKey(worker.Queue, strconv.Itoa(currentIndexToRead))

		jobToProcess := getValueByKey(key)

		if jobToProcess.IsProcessed {
		}

		if !jobToProcess.IsNotFound {

			if jobToProcess.Name != nil {
				jobName := *jobToProcess.Name
				jobData := *jobToProcess.Data

				log.Printf("Dispatching... %s", jobName)

				handler, ok := worker.Handlers[jobName]

				if !ok {
					log.Printf("No handler registered for job '%s' on queue '%s'", jobName, worker.Queue)
					setJobProcessed(worker.Queue, currentIndexToRead)
					continue
				}

				attemptsStartedForJob := getAttemptsCount(worker.Queue, currentIndexToRead, "ats")
				attemptsStartedForJob++
				setAttemptsCount(worker.Queue, currentIndexToRead, attemptsStartedForJob, "ats")

				var data map[string]any
				if err := json.Unmarshal([]byte(jobData), &data); err != nil {
					continue
				}

				setJobProcessed(worker.Queue, currentIndexToRead)

				returnValue, err := handler(ctx, data)

				if err != nil {
					fmt.Println(err)
				}

				finishedOn := time.Now().UnixMilli()
				setJobFinished(worker.Queue, currentIndexToRead, finishedOn, returnValue)
				addToJobCompletedQueue(worker.Queue, currentIndexToRead, finishedOn)

				attemptsMadeForJob := getAttemptsCount(worker.Queue, currentIndexToRead, "atm")
				attemptsMadeForJob++
				setAttemptsCount(worker.Queue, currentIndexToRead, attemptsMadeForJob, "atm")
			}

		}

		time.Sleep(10 * time.Millisecond)
	}
}
