package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sagar-pardhi/go-job-queue/config"
	"github.com/sagar-pardhi/go-job-queue/internal/database"
	"github.com/sagar-pardhi/go-job-queue/internal/jobs"
	redisclient "github.com/sagar-pardhi/go-job-queue/internal/redis"
	"github.com/sagar-pardhi/go-job-queue/internal/worker"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg.DatabaseURL())
	if err != nil {
		log.Fatal(err)
	}

	repo := jobs.NewRepository(db)

	redisClient := redisclient.NewRedis(
		cfg.RedisAddr,
	)

	log.Println("Worker started...")

	for {
		result, err := redisClient.BRPop(
			context.Background(),
			0,
			"queue:jobs",
		).Result()
		if err != nil {
			log.Println(err)
			continue
		}

		jobID := result[1]

		job, err := repo.GetByID(jobID)

		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("Started Processing: ", jobID)

		repo.UpdateStatus(jobID, jobs.StatusProcessing)

		err = worker.ProcessJob()

		if err != nil {
			retries := job.Retries + 1

			log.Printf(
				"Job %s failed. Retry %d/%d",
				jobID,
				retries,
				job.MaxRetries,
			)

			if retries < job.MaxRetries {
				repo.UpdateFailure(
					jobID,
					job.Retries+1,
					err.Error(),
				)

				delay := worker.RetryDelay(
					retries,
				)

				retryAt := time.Now().Add(delay).Unix()

				redisClient.ZAdd(
					context.Background(),
					"queue:delayed",
					redis.Z{
						Score:  float64(retryAt),
						Member: jobID,
					},
				)

				log.Printf("Job %s scheduled for retry in %v", jobID, delay)
			} else {
				repo.UpdateFailure(
					jobID,
					retries,
					err.Error(),
				)

				repo.UpdateStatus(
					jobID,
					jobs.StatusFailed,
				)

				log.Println("Job permanently failed: ", jobID)
			}
			continue
		}

		repo.ClearError(jobID)

		repo.UpdateStatus(
			jobID,
			jobs.StatusCompleted,
		)

		log.Println("Job Completed: ", jobID)
	}
}
