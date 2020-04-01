package telegram

import (
	"errors"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/yi-jiayu/ted"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"bot_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"telegram_bot": dataSourceTelegramBot(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"telegram_bot_webhook":  resourceTelegramBotWebhook(),
			"telegram_bot_commands": resourceTelegramBotCommands(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var token string
	if s, ok := d.Get("bot_token").(string); ok {
		token = s
	}
	if token == "" {
		token = os.Getenv("TELEGRAM_BOT_TOKEN")
	}
	if token == "" {
		return nil, errors.New("either bot_token or the environment variable TELEGRAM_BOT_TOKEN should be set")
	}
	bot := ted.Bot{
		Token:      token,
		HTTPClient: http.DefaultClient,
	}
	return bot, nil
}
