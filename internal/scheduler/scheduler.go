package scheduler

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Scheduler struct {
	redis *redis.Client
}

func New(redis *redis.Client) *Scheduler {
	return &Scheduler{
		redis: redis,
	}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(
		1 * time.Second,
	)

	defer ticker.Stop()

	for range ticker.C {
		now := strconv.FormatInt(
			time.Now().Unix(),
			10,
		)

		jobs, err := s.redis.ZRangeArgs(
			context.Background(),
			redis.ZRangeArgs{
				Key:     "queue:delayed",
				Start:   "-inf",
				Stop:    now,
				ByScore: true,
			},
		).Result()

		if err != nil {
			log.Println(err)
			continue
		}

		for _, jobID := range jobs {
			if err := s.redis.LPush(
				context.Background(),
				"queue:jobs",
				jobID,
			).Err(); err != nil {
				log.Println("failed to push job:", err)
				continue
			}

			if err := s.redis.ZRem(
				context.Background(),
				"queue:delayed",
				jobID,
			).Err(); err != nil {
				log.Println("failed to remove delayed job:", err)
			}

			log.Printf(
				"Moved delayed job %s to active queue",
				jobID,
			)
		}
	}
}
