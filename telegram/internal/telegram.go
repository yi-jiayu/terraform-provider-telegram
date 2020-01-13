package internal

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var retryableRegexp = regexp.MustCompile(`Too Many Requests: retry after (\d+)`)

func isRetryable(err error) (bool, int) {
	var telegramError tgbotapi.Error
	if ok := errors.As(err, &telegramError); ok && telegramError.RetryAfter > 0 {
		return true, telegramError.RetryAfter
	}
	// FIXME: remove after https://github.com/go-telegram-bot-api/telegram-bot-api/pull/300 is merged
	if matches := retryableRegexp.FindStringSubmatch(err.Error()); len(matches) > 0 {
		retryAfter, err := strconv.Atoi(matches[1])
		if err != nil {
			return false, 0
		}
		return true, retryAfter
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
