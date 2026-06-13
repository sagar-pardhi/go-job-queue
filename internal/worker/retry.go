package worker

import "time"

func RetryDelay(retries int) time.Duration {
	switch retries {
	case 1:
		return 5 * time.Second
	case 2:
		return 15 * time.Second
	default:
		return 30 * time.Second
	}
}
