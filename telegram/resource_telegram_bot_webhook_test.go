package telegram

import (
	"errors"
	"fmt"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceTelegramBotWebhook(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceTelegramBotWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
resource "telegram_bot_webhook" "example" {
  url = "https://www.example.com/webhook"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("telegram_bot_webhook.example", "url", "https://www.example.com/webhook"),
					testAccResourceTelegramBotWebhook("telegram_bot_webhook.example"),
				),
			},
			{
				Config: `
resource "telegram_bot_webhook" "example" {
  url = "https://www.example.com/newWebhook"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("telegram_bot_webhook.example", "url", "https://www.example.com/newWebhook"),
					testAccResourceTelegramBotWebhook("telegram_bot_webhook.example"),
				),
			},
			{
				PreConfig: func() {
					botAPI := testAccProvider.Meta().(*tgbotapi.BotAPI)
					_, err := botAPI.RemoveWebhook()
					if err != nil {
						t.Fatalf("error removing webhook: %s", err)
					}
				},
				Config: `
resource "telegram_bot_webhook" "example" {
  url = "https://www.example.com/newWebhook"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("telegram_bot_webhook.example", "url", "https://www.example.com/newWebhook"),
					testAccResourceTelegramBotWebhook("telegram_bot_webhook.example"),
				),
			},
		},
	})
}

func testAccResourceTelegramBotWebhook(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}
		botAPI := testAccProvider.Meta().(*tgbotapi.BotAPI)
		info, err := botAPI.GetWebhookInfo()
		if err != nil {
			return err
		}
		if got, want := info.URL, rs.Primary.Attributes["url"]; got != want {
			return fmt.Errorf("wanted webhook to be set to %s, got %s", want, got)
		}
		return nil
	}
}

func testAccResourceTelegramBotWebhookDestroy(*terraform.State) error {
	botAPI := testAccProvider.Meta().(*tgbotapi.BotAPI)
	info, err := botAPI.GetWebhookInfo()
	if err != nil {
		return err
	}
	if got, want := info.URL, ""; got != want {
		return errors.New("webhook was not destroyed")
	}
	return nil
}
