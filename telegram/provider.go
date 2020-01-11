package telegram

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"bot_token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
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
	s := d.Get("bot_token").(string)
	botAPI, err := tgbotapi.NewBotAPI(s)
	if err != nil {
		return nil, err
	}
	return BotAPI{botAPI}, nil
}
