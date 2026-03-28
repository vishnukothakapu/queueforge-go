package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"jobQueue-go/internal/model"
	"jobQueue-go/internal/queue"
	"jobQueue-go/internal/service"
	"jobQueue-go/pkg/db"
)

const WorkerCount = 5

type Job struct {
	ID         string
	Type       string
	Data       string
	Retries    int
	MaxRetries int
}

func process(job Job) error {
	fmt.Println("Processing job:", job.ID)
	service.UpdateJobStatus(job.ID, "processing")

	if time.Now().Unix()%2 == 0 {
		return fmt.Errorf("random failure")
	}
	time.Sleep(2 * time.Second)

	service.UpdateJobStatus(job.ID, "completed")

	fmt.Println("Completed job:", job.ID)
	return nil
}

func worker(ctx context.Context, id int, jobs <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d shutting down...\n", id)
			return
		case job := <-jobs:
			fmt.Printf("Worker %d picked job %s\n", id, job.ID)
			err := process(job)
			if err != nil {
				fmt.Println("Job failed", job.ID)
				job.Retries++

				if job.Retries >= job.MaxRetries {
					service.UpdateJobStatus(job.ID, "failed")
					fmt.Println("Job permanently failed:", job.ID)
				} else {
					service.UpdateJobStatus(job.ID, "retrying")

					time.Sleep(3 * time.Second)
					queue.Enqueue(model.Job{
						ID:         job.ID,
						Type:       job.Type,
						Data:       job.Data,
						Status:     "queued",
						Retries:    job.Retries,
						MaxRetries: job.MaxRetries,
					})
				}
			}
		}
	}
}

func main() {
	db.Init()
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	jobs := make(chan Job, 100)

	var wg sync.WaitGroup
	for i := 1; i <= WorkerCount; i++ {
		go worker(ctx, i, jobs, &wg)
	}
	go func() {
		<-sigChan
		fmt.Print("Shutdown signal received...")

		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping job fetcher...")
			wg.Wait()
			fmt.Println("All workers stopped.")
			return

		default:
			job, err := queue.Dequeue()
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}
			jobs <- Job{
				ID:         job.ID,
				Type:       job.Type,
				Data:       job.Data,
				Retries:    job.Retries,
				MaxRetries: job.MaxRetries,
			}

		}
	}
}
