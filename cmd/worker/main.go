package main

import (
	"context"
	"log"
	"time"

	"github.com/sagar-pardhi/go-job-queue/internal/database"
	"github.com/sagar-pardhi/go-job-queue/internal/jobs"
	redisclient "github.com/sagar-pardhi/go-job-queue/internal/redis"
)

func main() {
	connString := database.BuildConnectionString(
		"localhost",
		"5433",
		"admin",
		"admin",
		"jobqueue",
	)

	db, err := database.NewPostgres(connString)
	if err != nil {
		log.Fatal(err)
	}

	repo := jobs.NewRepository(db)

	redisClient := redisclient.NewRedis(
		"localhost:6379",
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

		repo.UpdateStatus(jobID, "processing")

		time.Sleep(5 * time.Second)

		repo.UpdateStatus(jobID, "completed")

		log.Println("Completed: ", jobID)
	}
}
