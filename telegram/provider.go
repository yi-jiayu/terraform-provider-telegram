package telegram

import (
	"errors"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

type BotAPI struct {
	*tgbotapi.BotAPI
}

func (bot BotAPI) ID() string {
	return strconv.Itoa(bot.Self.ID)
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
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return BotAPI{botAPI}, nil
}
