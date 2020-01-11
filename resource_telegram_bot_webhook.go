package main

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//go:generate mockgen -source=resource_telegram_bot_webhook.go -destination=mocks/resource_telegram_bot_webhook.go -package=mocks

type webhookSetter interface {
	ID() string
	SetWebhook(config tgbotapi.WebhookConfig) (tgbotapi.APIResponse, error)
}

type webhookGetter interface {
	GetWebhookInfo() (tgbotapi.WebhookInfo, error)
}

type webhookRemover interface {
	RemoveWebhook() (tgbotapi.APIResponse, error)
}

func resourceTelegramBotWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceTelegramBotWebhookCreate,
		Read:   resourceTelegramBotWebhookRead,
		Update: resourceTelegramBotWebhookUpdate,
		Delete: resourceTelegramBotWebhookDelete,
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_connections": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceTelegramBotWebhookCreate(d *schema.ResourceData, m interface{}) error {
	setter := m.(webhookSetter)
	url := d.Get("url").(string)
	result, err := setter.SetWebhook(tgbotapi.NewWebhook(url))
	if err != nil {
		return err
	}
	if !result.Ok {
		return errors.New(result.Description)
	}
	d.SetId(setter.ID())
	return nil
}

func resourceTelegramBotWebhookRead(d *schema.ResourceData, m interface{}) error {
	getter := m.(webhookGetter)
	info, err := getter.GetWebhookInfo()
	if err != nil {
		return err
	}
	url := info.URL
	if url == "" {
		d.SetId("")
		return nil
	}
	return d.Set("url", url)
}

func resourceTelegramBotWebhookUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceTelegramBotWebhookCreate(d, m)
}

func resourceTelegramBotWebhookDelete(d *schema.ResourceData, m interface{}) error {
	remover := m.(webhookRemover)
	result, err := remover.RemoveWebhook()
	if err != nil {
		return err
	}
	if !result.Ok {
		return errors.New(result.Description)
	}
	d.SetId("")
	return nil
}
