package telegram

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/yi-jiayu/ted"
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
	bot := m.(ted.Bot)
	me, err := bot.GetMe()
	if err != nil {
		return fmt.Errorf("getMe error: %w", err)
	}
	d.SetId(strconv.FormatInt(me.ID, 10))
	if err := d.Set("user_id", me.ID); err != nil {
		return err
	}
	if err := d.Set("name", me.FirstName); err != nil {
		return err
	}
	if err := d.Set("username", me.Username); err != nil {
		return err
	}
	return nil
}
