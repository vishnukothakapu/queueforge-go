package queue

import (
	"encoding/json"
	"jobQueue-go/internal/model"
	"jobQueue-go/pkg/redis"
)

const QueueName = "job_queue"

func Enqueue(job model.Job) error {
	client := redis.NewClient()

	data, _ := json.Marshal(job)

	return client.RPush(redis.Ctx, QueueName, data).Err()
}

func Dequeue() (model.Job, error) {
	client := redis.NewClient()

	result, err := client.LPop(redis.Ctx, QueueName).Result()
	if err != nil {
		return model.Job{}, err
	}

	var job model.Job
	json.Unmarshal([]byte(result), &job)

	return job, nil

}
