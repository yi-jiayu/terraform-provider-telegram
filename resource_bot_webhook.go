package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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
	url := d.Get("url").(string)
	d.SetId(url)
	return resourceBotWebhookRead(d, m)
}

func resourceBotWebhookRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceBotWebhookUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceBotWebhookRead(d, m)
}

func resourceBotWebhookDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
