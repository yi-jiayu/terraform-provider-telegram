package telegram

import (
	"errors"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/yi-jiayu/terraform-provider-telegram/telegram/internal"
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
			"telegram_bot_webhook": resourceTelegramBotWebhook(),
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
	var botAPI *tgbotapi.BotAPI
	err := internal.Retry(3, func() error {
		var err error
		botAPI, err = tgbotapi.NewBotAPI(token)
		if err != nil {
			return fmt.Errorf("error creating Telegram API client: %w", err)
		}
		return nil
	})
	return botAPI, err
}
