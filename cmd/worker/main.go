package main

import (
	"context"
	"log"

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

	service := worker.NewService(
		repo,
		redisClient,
	)

	jobsChan := make(chan string, 100)

	workerCount := 10

	for i := 1; i <= workerCount; i++ {
		go worker.StartWorker(
			i,
			service,
			jobsChan,
		)
	}

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

		jobsChan <- jobID
	}
}
