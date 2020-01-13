package telegram

import (
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Test_retry(t *testing.T) {
	f := func() error {
		return tgbotapi.Error{
			Message: "error",
			ResponseParameters: tgbotapi.ResponseParameters{
				RetryAfter: 1,
			},
		}
	}
	err := retry(0, 3, f)
	if err == nil {
		t.Fatal("wanted err to be non-nil")
	}
	if got, want := err.Error(), "retried 3 times: "; !strings.Contains(got, want) {
		t.Fatalf("wanted %s to contain %s", got, want)
	}
}
