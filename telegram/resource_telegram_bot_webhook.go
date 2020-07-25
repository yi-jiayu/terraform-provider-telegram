package telegram

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/yi-jiayu/ted"

	"github.com/yi-jiayu/terraform-provider-telegram/telegram/internal"
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
			"max_connections": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  40,
			},
			"allowed_updates": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceTelegramBotWebhookCreate(d *schema.ResourceData, m interface{}) error {
	bot := m.(ted.Bot)
	url := d.Get("url").(string)
	maxConnections := d.Get("max_connections").(int)
	setWebhook := ted.SetWebhookRequest{
		URL:            url,
		MaxConnections: maxConnections,
		AllowedUpdates: []string{},
	}
	if v, ok := d.GetOk("allowed_updates"); ok {
		vs := v.(*schema.Set).List()
		allowedUpdates := make([]string, len(vs))
		for i, updateType := range vs {
			allowedUpdates[i] = updateType.(string)
		}
		setWebhook.AllowedUpdates = allowedUpdates
	}
	err := internal.Retry(3, func() error {
		_, err := bot.Do(setWebhook)
		if err != nil {
			return fmt.Errorf("setWebhook error: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId("webhook")
	return resourceTelegramBotWebhookRead(d, m)
}

func resourceTelegramBotWebhookRead(d *schema.ResourceData, m interface{}) error {
	bot := m.(ted.Bot)
	var info ted.WebhookInfo
	err := internal.Retry(3, func() error {
		res, err := bot.Do(ted.GetWebhookInfoRequest{})
		if err != nil {
			return fmt.Errorf("getWebhookInfo error: %w", err)
		}
		err = json.Unmarshal(res.Result, &info)
		if err != nil {
			return err
		}
		return nil
	})
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
	if err := d.Set("max_connections", info.MaxConnections); err != nil {
		return err
	}
	if len(info.AllowedUpdates) > 0 {
		allowedUpdates := schema.NewSet(schema.HashString, nil)
		for _, updateType := range info.AllowedUpdates {
			allowedUpdates.Add(updateType)
		}
		if err := d.Set("allowed_updates", info.AllowedUpdates); err != nil {
			return err
		}
	}
	return nil
}

func resourceTelegramBotWebhookUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceTelegramBotWebhookCreate(d, m)
}

func resourceTelegramBotWebhookDelete(d *schema.ResourceData, m interface{}) error {
	bot := m.(ted.Bot)
	removeWebhook := ted.SetWebhookRequest{}
	err := internal.Retry(3, func() error {
		_, err := bot.Do(removeWebhook)
		if err != nil {
			return fmt.Errorf("removeWebhook error: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
