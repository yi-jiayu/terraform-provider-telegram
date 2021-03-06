package internal

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/yi-jiayu/ted"
)

func Test_retry(t *testing.T) {
	f := func() error {
		return ted.Response{
			Parameters: &ted.ResponseParameters{
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

func Test_isRetryable(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		wantOk         bool
		wantRetryAfter int
	}{
		{
			name: "429 error from Telegram",
			err: ted.Response{
				Parameters: &ted.ResponseParameters{
					RetryAfter: 1,
				},
			},
			wantOk:         true,
			wantRetryAfter: 1,
		},
		{
			name:           "wrapped 429 error from Telegram",
			err:            fmt.Errorf("%w", ted.Response{Parameters: &ted.ResponseParameters{RetryAfter: 1}}),
			wantOk:         true,
			wantRetryAfter: 1,
		},
		{
			name:           "not retryable",
			err:            errors.New("random error"),
			wantOk:         false,
			wantRetryAfter: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, retryAfter := isRetryable(tt.err)
			if ok != tt.wantOk {
				t.Fatalf("want %t, got %t", tt.wantOk, ok)
			}
			if retryAfter != tt.wantRetryAfter {
				t.Fatalf("want %d, got %d", tt.wantRetryAfter, retryAfter)
			}
		})
	}
}
