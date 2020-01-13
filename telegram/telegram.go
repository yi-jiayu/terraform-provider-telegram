package telegram

import (
	"errors"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func retry(numRetries, maxRetries int, f func() error) error {
	err := f()
	if err != nil {
		var telegramError tgbotapi.Error
		if ok := errors.As(err, &telegramError); ok && telegramError.RetryAfter > 0 && numRetries < maxRetries {
			time.Sleep(time.Duration(telegramError.RetryAfter) * time.Second)
			return retry(numRetries+1, maxRetries, f)
		}
		return fmt.Errorf("retried %d times: %w", numRetries, err)
	}
	return nil
}
