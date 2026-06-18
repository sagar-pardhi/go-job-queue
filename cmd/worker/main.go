package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sagar-pardhi/go-job-queue/config"
	"github.com/sagar-pardhi/go-job-queue/internal/database"
	"github.com/sagar-pardhi/go-job-queue/internal/jobs"
	redisclient "github.com/sagar-pardhi/go-job-queue/internal/redis"
	"github.com/sagar-pardhi/go-job-queue/internal/worker"
)

var wg sync.WaitGroup

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	defer cancel()

	sigChan := make(
		chan os.Signal,
		1,
	)

	signal.Notify(
		sigChan,
		os.Interrupt,
		syscall.SIGTERM,
	)

	go func() {
		sig := <-sigChan

		log.Printf("Received signal: %v", sig)

		cancel()
	}()

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
		wg.Add(1)
		go worker.StartWorker(
			ctx,
			i,
			service,
			jobsChan,
			&wg,
		)
	}

	for {

		select {

		case <-ctx.Done():

			log.Println(
				"Stopping job fetcher...",
			)

			close(jobsChan)

			wg.Wait()

			log.Println(
				"All workers stopped",
			)

			return

		default:

			result, err := redisClient.BRPop(
				context.Background(),
				1*time.Second,
				"queue:jobs",
			).Result()

			if err != nil {

				if err == redis.Nil {
					continue
				}

				log.Println(err)

				continue
			}

			jobID := result[1]

			jobsChan <- jobID
		}
	}
}
