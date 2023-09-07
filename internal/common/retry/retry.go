package retry

import (
	"context"
	"errors"
	"net/http"
	"time"
)

func Retry(ctx context.Context, fn func() (*http.Response, error), retries int) (*http.Response, error) {
	var errs []error
	for i := 0; i < retries; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			r, err := fn()
			if err == nil {
				return r, nil
			}
			errs = append(errs, err)
			interval := time.Duration(i+1) * time.Second
			time.Sleep(interval)
		}
	}
	err := errors.New("max retries reached")
	errs = append(errs, err)
	return nil, errors.Join(errs...)
}
