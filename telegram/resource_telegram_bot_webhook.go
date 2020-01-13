package telegram

import (
	"errors"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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
			"has_custom_certificate": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceTelegramBotWebhookCreate(d *schema.ResourceData, m interface{}) error {
	botAPI := m.(*tgbotapi.BotAPI)
	url := d.Get("url").(string)
	config := tgbotapi.NewWebhook(url)
	if cert, ok := d.Get("certificate").(string); ok && cert != "" {
		config.Certificate = tgbotapi.FileBytes{Bytes: []byte(cert)}
	}
	result, err := botAPI.SetWebhook(config)
	if err != nil {
		return err
	}
	if !result.Ok {
		return errors.New(result.Description)
	}
	d.SetId(strconv.Itoa(botAPI.Self.ID))
	return resourceTelegramBotWebhookRead(d, m)
}

func resourceTelegramBotWebhookRead(d *schema.ResourceData, m interface{}) error {
	botAPI := m.(*tgbotapi.BotAPI)
	info, err := botAPI.GetWebhookInfo()
	if err != nil {
		return err
	}
	url := info.URL
	if url == "" {
		d.SetId("")
		return nil
	}
	if err := d.Set("url", url); err != nil {
		return err
	}
	if err := d.Set("has_custom_certificate", info.HasCustomCertificate); err != nil {
		return err
	}
	return nil
}

func resourceTelegramBotWebhookUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceTelegramBotWebhookCreate(d, m)
}

func resourceTelegramBotWebhookDelete(d *schema.ResourceData, m interface{}) error {
	botAPI := m.(*tgbotapi.BotAPI)
	result, err := botAPI.RemoveWebhook()
	if err != nil {
		return err
	}
	if !result.Ok {
		return errors.New(result.Description)
	}
	d.SetId("")
	return nil
}
