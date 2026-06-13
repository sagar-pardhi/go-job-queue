package worker

import "log"

func StartWorker(id int, service *Service, jobs <-chan string) {
	for jobID := range jobs {
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
				jobID,
			)
		}
	}
}
