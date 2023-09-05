package retry

import (
	"errors"
	"time"
)

func Retry(fn func() error, retries int) error {
	for i := 0; i < retries; i++ {
		err := fn()
		if err == nil {
			return nil
		}
		interval := time.Duration(i+1) * time.Second
		time.Sleep(interval)
	}
	return errors.New("max retries reached")

}
