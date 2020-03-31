package internal

import (
	"errors"
	"fmt"
	"time"

	"github.com/yi-jiayu/ted"
)

// isRetryable returns a boolean indicating whether an error is retryable.
// If so, the second return value indicates how many seconds to wait before retrying.
func isRetryable(err error) (bool, int) {
	var telegramError ted.Response
	if ok := errors.As(err, &telegramError); ok && telegramError.Parameters.RetryAfter > 0 {
		return true, telegramError.Parameters.RetryAfter
	}
	return false, 0
}

func retry(numRetries, maxRetries int, f func() error) error {
	err := f()
	if err != nil {
		if ok, retryAfter := isRetryable(err); ok && numRetries < maxRetries {
			time.Sleep(time.Duration(retryAfter) * time.Second)
			return retry(numRetries+1, maxRetries, f)
		}
		return fmt.Errorf("retried %d times: %w", numRetries, err)
	}
	return nil
}

func Retry(maxRetries int, f func() error) error {
	return retry(0, maxRetries, f)
}
