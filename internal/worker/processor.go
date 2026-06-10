package worker

import (
	"errors"
	"math/rand"
)

func ProcessJob() error {
	if rand.Intn(10) < 3 {
		return errors.New("simulated failure")
	}
	return nil
}
