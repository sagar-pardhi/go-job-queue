package worker

import (
	"context"
	"log"
	"sync"
)

func StartWorker(ctx context.Context, id int, service *Service, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {

		select {

		case <-ctx.Done():

			log.Printf(
				"[Worker %d] shutting down",
				id,
			)

			return

		case jobID, ok := <-jobs:

			if !ok {
				return
			}

			log.Printf(
				"[Worker %d] Processing %s",
				id,
				jobID,
			)

			err := service.Process(
				jobID,
			)

			if err != nil {
				log.Printf(
					"[Worker %d] Error: %v",
					id,
					err,
				)
			}
		}
	}
}
