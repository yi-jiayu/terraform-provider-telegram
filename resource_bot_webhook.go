package main

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type WebhookSetter interface {
	ID() string
	SetWebhook(config tgbotapi.WebhookConfig) (tgbotapi.APIResponse, error)
}

type WebhookGetter interface {
	GetWebhookInfo() (tgbotapi.WebhookInfo, error)
}

type WebhookRemover interface {
	RemoveWebhook() (tgbotapi.APIResponse, error)
}

func resourceBotWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceBotWebhookCreate,
		Read:   resourceBotWebhookRead,
		Update: resourceBotWebhookUpdate,
		Delete: resourceBotWebhookDelete,
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceBotWebhookCreate(d *schema.ResourceData, m interface{}) error {
	setter := m.(WebhookSetter)
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

func resourceBotWebhookRead(d *schema.ResourceData, m interface{}) error {
	getter := m.(WebhookGetter)
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

func resourceBotWebhookUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceBotWebhookCreate(d, m)
}

func resourceBotWebhookDelete(d *schema.ResourceData, m interface{}) error {
	remover := m.(WebhookRemover)
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
