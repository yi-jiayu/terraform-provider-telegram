package main

import (
	"errors"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/yi-jiayu/terraform-provider-telegram/mocks"
)

func Test_resourceTelegramBotWebhookCreate(t *testing.T) {
	tests := []struct {
		name     string
		state    map[string]interface{}
		response tgbotapi.APIResponse
		err      error
		assert   func(t *testing.T, d *schema.ResourceData, err error)
	}{
		{
			name:     "webhook created",
			response: tgbotapi.APIResponse{Ok: true},
			assert: func(t *testing.T, d *schema.ResourceData, err error) {
				if err != nil {
					t.Fatal("wanted err to be nil")
				}
				want := "1234"
				if got := d.Id(); got != want {
					t.Fatalf("wanted id to be %s, got %s", want, got)
				}
			},
		},
		{
			name:     "error setting webhook",
			response: tgbotapi.APIResponse{},
			err:      errors.New("error"),
			assert: func(t *testing.T, d *schema.ResourceData, err error) {
				if err == nil {
					t.Fatal("wanted err to be non-nil")
				}
				want := "error"
				if got := err.Error(); got != want {
					t.Fatalf("wanted error to be %s, got %s", want, got)
				}
			},
		},
		{
			name:     "response not ok",
			response: tgbotapi.APIResponse{Ok: false, Description: "description"},
			assert: func(t *testing.T, d *schema.ResourceData, err error) {
				if err == nil {
					t.Fatal("wanted err to be non-nil")
				}
				want := "description"
				if got := err.Error(); got != want {
					t.Fatalf("wanted error to be %s, got %s", want, got)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockwebhookSetter(ctrl)
			m.EXPECT().SetWebhook(gomock.Any()).Return(tt.response, tt.err)
			m.EXPECT().ID().Return("1234").AnyTimes()

			d := schema.TestResourceDataRaw(t, resourceTelegramBotWebhook().Schema, tt.state)
			err := resourceTelegramBotWebhookCreate(d, m)
			tt.assert(t, d, err)
		})
	}
}

func Test_resourceTelegramBotWebhookRead(t *testing.T) {
	tests := []struct {
		name   string
		state  map[string]interface{}
		info   tgbotapi.WebhookInfo
		err    error
		assert func(t *testing.T, d *schema.ResourceData, err error)
	}{
		{
			name: "read success",
			state: map[string]interface{}{
				"url": "https://www.example.com/oldWebhook",
			},
			info: tgbotapi.WebhookInfo{URL: "https://www.example.com/newWebhook"},
			err:  nil,
			assert: func(t *testing.T, d *schema.ResourceData, err error) {
				if err != nil {
					t.Fatal("wanted err to be nil")
				}
				want := "https://www.example.com/newWebhook"
				if got := d.Get("url").(string); got != want {
					t.Fatalf("wanted url to be %s, got %s", want, got)
				}
			},
		},
		{
			name: "webhook removed",
			state: map[string]interface{}{
				"id": "123",
			},
			info: tgbotapi.WebhookInfo{URL: ""},
			err:  nil,
			assert: func(t *testing.T, d *schema.ResourceData, err error) {
				if err != nil {
					t.Fatal("wanted err to be nil")
				}
				if d.Id() != "" {
					t.Fatal("wanted id to be set to empty string")
				}
			},
		},
		{
			name:  "error getting webhook response",
			state: nil,
			info:  tgbotapi.WebhookInfo{},
			err:   errors.New("error"),
			assert: func(t *testing.T, d *schema.ResourceData, err error) {
				if err == nil {
					t.Fatal("wanted err to be non-nil")
				}
				want := "error"
				if got := err.Error(); got != want {
					t.Fatalf("wanted err to be %s, got %s", want, got)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockwebhookGetter(ctrl)
			m.EXPECT().GetWebhookInfo().Return(tt.info, tt.err)

			d := schema.TestResourceDataRaw(t, resourceTelegramBotWebhook().Schema, tt.state)
			err := resourceTelegramBotWebhookRead(d, m)
			tt.assert(t, d, err)
		})
	}
}

func Test_resourceTelegramBotWebhookDelete(t *testing.T) {
	tests := []struct {
		name     string
		state    map[string]interface{}
		response tgbotapi.APIResponse
		err      error
		assert   func(t *testing.T, d *schema.ResourceData, err error)
	}{
		{
			name:     "webhook deleted",
			response: tgbotapi.APIResponse{Ok: true},
			assert: func(t *testing.T, d *schema.ResourceData, err error) {
				if err != nil {
					t.Fatal("wanted err to be nil")
				}
			},
		},
		{
			name:     "error deleting webhook",
			response: tgbotapi.APIResponse{},
			err:      errors.New("error"),
			assert: func(t *testing.T, d *schema.ResourceData, err error) {
				if err == nil {
					t.Fatal("wanted err to be non-nil")
				}
				want := "error"
				if got := err.Error(); got != want {
					t.Fatalf("wanted error to be %s, got %s", want, got)
				}
			},
		},
		{
			name:     "response not ok",
			response: tgbotapi.APIResponse{Ok: false, Description: "description"},
			assert: func(t *testing.T, d *schema.ResourceData, err error) {
				if err == nil {
					t.Fatal("wanted err to be non-nil")
				}
				want := "description"
				if got := err.Error(); got != want {
					t.Fatalf("wanted error to be %s, got %s", want, got)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockwebhookRemover(ctrl)
			m.EXPECT().RemoveWebhook().Return(tt.response, tt.err)

			d := schema.TestResourceDataRaw(t, resourceTelegramBotWebhook().Schema, tt.state)
			err := resourceTelegramBotWebhookDelete(d, m)
			tt.assert(t, d, err)
		})
	}
}
