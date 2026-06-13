package worker

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sagar-pardhi/go-job-queue/internal/jobs"
)

type Service struct {
	repo  *jobs.Repository
	redis *redis.Client
}

func NewService(repo *jobs.Repository, redis *redis.Client) *Service {
	return &Service{
		repo:  repo,
		redis: redis,
	}
}

func (s *Service) Process(jobID string) error {
	job, err := s.repo.GetByID(jobID)

	if err != nil {
		log.Printf(
			"Failed to load job %s: %v",
			jobID,
			err,
		)

		return err
	}

	log.Printf(
		"Loaded job: id=%s retries=%d maxRetries=%d",
		job.ID,
		job.Retries,
		job.MaxRetries,
	)

	err = s.repo.UpdateStatus(
		jobID,
		jobs.StatusProcessing,
	)

	if err != nil {
		return err
	}

	err = ProcessJob()

	if err != nil {
		retries := job.Retries + 1

		log.Printf(
			"Updating retry count: job=%s retries=%d",
			jobID,
			retries,
		)

		if retries <= job.MaxRetries {
			log.Printf(
				"DB before update: retries=%d",
				job.Retries,
			)

			err := s.repo.UpdateFailure(
				jobID,
				retries,
				err.Error(),
			)

			updatedJob, getErr := s.repo.GetByID(jobID)

			if getErr != nil {
				log.Printf("Failed to reload job: %v", getErr)
			} else {
				log.Printf(
					"After update: retries=%d status=%s",
					updatedJob.Retries,
					updatedJob.Status,
				)
			}

			if err != nil {
				return err
			}

			err = s.repo.UpdateStatus(
				jobID,
				jobs.StatusRetrying,
			)

			if err != nil {
				return err
			}

			delay := RetryDelay(
				retries,
			)

			retryAt := time.Now().
				Add(delay).
				Unix()

			err = s.redis.ZAdd(
				context.Background(),
				"queue:delayed",
				redis.Z{
					Score:  float64(retryAt),
					Member: jobID,
				},
			).Err()

			if err != nil {
				return err
			}

			log.Printf(
				"Job %s scheduled in %v",
				jobID,
				delay,
			)

			return nil
		}

		log.Printf(
			"Updating retry count: job=%s retries=%d",
			jobID,
			retries,
		)

		err = s.repo.UpdateFailure(
			jobID,
			retries,
			err.Error(),
		)

		updatedJob, getErr := s.repo.GetByID(jobID)

		if getErr != nil {
			log.Printf("Failed to reload job: %v", getErr)
		} else {
			log.Printf(
				"After update: retries=%d status=%s",
				updatedJob.Retries,
				updatedJob.Status,
			)
		}

		if err != nil {
			return err
		}

		log.Printf("Job %s failed", jobID)

		return s.repo.UpdateStatus(
			jobID,
			jobs.StatusFailed,
		)
	}

	err = s.repo.ClearError(
		jobID,
	)

	if err != nil {
		return err
	}

	log.Printf("Job %s completed", jobID)

	return s.repo.UpdateStatus(
		jobID,
		jobs.StatusCompleted,
	)
}
