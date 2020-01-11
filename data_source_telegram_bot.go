package main

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTelegramBot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTelegramBotRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTelegramBotRead(d *schema.ResourceData, m interface{}) error {
	botAPI := m.(BotAPI)
	self := botAPI.Self
	d.SetId(strconv.Itoa(self.ID))
	if err := d.Set("user_id", self.ID); err != nil {
		return err
	}
	if err := d.Set("name", self.FirstName); err != nil {
		return err
	}
	if err := d.Set("username", self.UserName); err != nil {
		return err
	}
	return nil
}
