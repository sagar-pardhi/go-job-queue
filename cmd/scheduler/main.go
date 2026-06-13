package main

import (
	"log"

	"github.com/sagar-pardhi/go-job-queue/config"
	redisclient "github.com/sagar-pardhi/go-job-queue/internal/redis"
	"github.com/sagar-pardhi/go-job-queue/internal/scheduler"
)

func main() {
	cfg := config.Load()

	redisclient := redisclient.NewRedis(
		cfg.RedisAddr,
	)

	log.Println("Scheduler started")

	s := scheduler.New(
		redisclient,
	)

	s.Start()
}
