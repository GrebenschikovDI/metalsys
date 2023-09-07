package retry

import (
	"context"
	"errors"
	"time"
)

func Retry(ctx context.Context, fn func() error, retries int) error {
	var errs []error
	for i := 0; i < retries; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := fn()
			if err == nil {
				return nil
			}
			errs = append(errs, err)
			interval := time.Duration(i+1) * time.Second
			time.Sleep(interval)
		}
	}
	err := errors.New("max retries reached")
	errs = append(errs, err)
	return errors.Join(errs...)
}
