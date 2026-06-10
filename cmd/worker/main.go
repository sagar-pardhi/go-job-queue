package main

import (
	"context"
	"log"
	"time"

	"github.com/sagar-pardhi/go-job-queue/config"
	"github.com/sagar-pardhi/go-job-queue/internal/database"
	"github.com/sagar-pardhi/go-job-queue/internal/jobs"
	redisclient "github.com/sagar-pardhi/go-job-queue/internal/redis"
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

		log.Println("Processing: ", jobID)

		repo.UpdateStatus(jobID, jobs.StatusProcessing)

		time.Sleep(5 * time.Second)

		repo.UpdateStatus(jobID, jobs.StatusCompleted)

		log.Println("Completed: ", jobID)
	}
}
