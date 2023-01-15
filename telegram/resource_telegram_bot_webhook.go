package telegram

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/yi-jiayu/ted"

	"github.com/yi-jiayu/terraform-provider-telegram/telegram/internal"
)

func resourceTelegramBotWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTelegramBotWebhookCreate,
		ReadContext:   resourceTelegramBotWebhookRead,
		UpdateContext: resourceTelegramBotWebhookUpdate,
		Delete:        resourceTelegramBotWebhookDelete,
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

func resourceTelegramBotWebhookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}
	d.SetId("webhook")
	return resourceTelegramBotWebhookRead(ctx, d, m)
}

func resourceTelegramBotWebhookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}
	url := info.URL
	if url == "" {
		d.SetId("")
		return nil
	}
	if err := d.Set("url", url); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("max_connections", info.MaxConnections); err != nil {
		return diag.FromErr(err)
	}
	if len(info.AllowedUpdates) > 0 {
		allowedUpdates := schema.NewSet(schema.HashString, nil)
		for _, updateType := range info.AllowedUpdates {
			allowedUpdates.Add(updateType)
		}
		if err := d.Set("allowed_updates", info.AllowedUpdates); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceTelegramBotWebhookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceTelegramBotWebhookCreate(ctx, d, m)
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
