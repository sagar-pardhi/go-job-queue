package worker

import (
	"errors"
	"math/rand"
	"time"
)

func ProcessJob() error {
	time.Sleep(5 * time.Second)

	if rand.Intn(10) < 3 {
		return errors.New("simulated failure")
	}
	return nil
}
