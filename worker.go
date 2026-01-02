package oncamq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Handlers map[string]func(ctx context.Context, data map[string]any) (any, error)

type Worker struct {
	redisClient *redis.Client
	Queue       string
	Handlers    Handlers
}

type GetKeyFromQueueResponse struct {
	Name        *string
	Data        *string
	IsProcessed bool
	IsNotFound  bool
}

func createKey(queue string, key string) string {
	return fmt.Sprintf("bull:%s:%s", queue, key)
}

func (w *Worker) getValueByKey(ctx context.Context, key string) GetKeyFromQueueResponse {

	job, err := w.redisClient.HGetAll(ctx, key).Result()

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

	value, err := w.redisClient.HGet(ctx, key, "data").Result()

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

func (w *Worker) setJobProcessed(ctx context.Context, queue string, key int) {
	fullKeyJob := fmt.Sprintf("bull:%s:%d", queue, key)

	w.redisClient.HSet(ctx, fullKeyJob,
		"processedOn", time.Now().UnixMilli(),
	)
}

func (w *Worker) setJobFinished(ctx context.Context, queue string, key int, finishedOn int64, returnValue interface{}) {
	fullKeyJob := fmt.Sprintf("bull:%s:%d", queue, key)

	w.redisClient.HSet(ctx, fullKeyJob,
		"finishedOn", finishedOn,
		"returnvalue", returnValueToString(returnValue),
	)
}

func (w *Worker) setAttemptsCount(ctx context.Context, queue string, key int, attempt int, attemptType string) {

	fullKeyJob := fmt.Sprintf("bull:%s:%d", queue, key)

	w.redisClient.HSet(ctx, fullKeyJob,
		attemptType, attempt,
	)
}

func (w *Worker) getLastIndexRead(ctx context.Context, queue string) int {

	fullKey := createLastReadKey(queue)

	value, err := w.redisClient.ZRange(ctx, fullKey, 0, -1).Result()

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

func (w *Worker) getAttemptsCount(ctx context.Context, queue string, key int, attemptType string) int {

	fullKeyJob := fmt.Sprintf("bull:%s:%d", queue, key)

	value, err := w.redisClient.HGet(ctx, fullKeyJob, attemptType).Result()

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

func (w *Worker) addToJobCompletedQueue(ctx context.Context, queue string, key int, finishedOn int64) {
	fullKeyJob := fmt.Sprintf("bull:%s:completed", queue)

	w.redisClient.ZAdd(ctx, fullKeyJob, redis.Z{
		Score:  float64(finishedOn),
		Member: key,
	})
}

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

func New(redisClient *redis.Client, queue string, handlers Handlers) *Worker {
	return &Worker{
		redisClient: redisClient,
		Queue:       queue,
		Handlers:    handlers,
	}
}

func (w *Worker) StartWorker(ctx context.Context) {
	log.Printf("Worker to consume %s (%s) started...", w.Queue, w.redisClient.Options().Addr)

	for {
		lastIndexRead := w.getLastIndexRead(ctx, w.Queue)

		currentIndexToRead := lastIndexRead + 1
		key := createKey(w.Queue, strconv.Itoa(currentIndexToRead))

		jobToProcess := w.getValueByKey(ctx, key)

		if jobToProcess.IsProcessed {
		}

		if !jobToProcess.IsNotFound {

			if jobToProcess.Name != nil {
				jobName := *jobToProcess.Name
				jobData := *jobToProcess.Data

				log.Printf("Dispatching job %s of %s (jobID: %d)", jobName, w.Queue, currentIndexToRead)

				handler, ok := w.Handlers[jobName]

				if !ok {
					log.Printf("No handler registered for job '%s' on queue '%s'", jobName, w.Queue)
					w.setJobProcessed(ctx, w.Queue, currentIndexToRead)
					continue
				}

				attemptsStartedForJob := w.getAttemptsCount(ctx, w.Queue, currentIndexToRead, "ats")
				attemptsStartedForJob++
				w.setAttemptsCount(ctx, w.Queue, currentIndexToRead, attemptsStartedForJob, "ats")

				var data map[string]any
				if err := json.Unmarshal([]byte(jobData), &data); err != nil {
					continue
				}

				w.setJobProcessed(ctx, w.Queue, currentIndexToRead)

				returnValue, err := handler(ctx, data)

				if err != nil {
					fmt.Println(err)
				}

				finishedOn := time.Now().UnixMilli()
				w.setJobFinished(ctx, w.Queue, currentIndexToRead, finishedOn, returnValue)
				w.addToJobCompletedQueue(ctx, w.Queue, currentIndexToRead, finishedOn)

				attemptsMadeForJob := w.getAttemptsCount(ctx, w.Queue, currentIndexToRead, "atm")
				attemptsMadeForJob++
				w.setAttemptsCount(ctx, w.Queue, currentIndexToRead, attemptsMadeForJob, "atm")
			}
		}
	}
}
