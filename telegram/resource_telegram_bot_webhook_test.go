package telegram

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/yi-jiayu/ted"
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
					bot := testAccProvider.Meta().(ted.Bot)
					_, err := bot.Do(ted.SetWebhookRequest{})
					if err != nil {
						t.Fatalf("error removing webhook: %s", err)
					}
				},
				Config: `
resource "telegram_bot_webhook" "example" {
  url = "https://www.example.com/newWebhook"
  allowed_updates = ["message", "inline_query"]
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("telegram_bot_webhook.example", "url", "https://www.example.com/newWebhook"),
					resource.TestCheckResourceAttr("telegram_bot_webhook.example", "allowed_updates.#", "2"),
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
		botAPI := testAccProvider.Meta().(ted.Bot)
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
	botAPI := testAccProvider.Meta().(ted.Bot)
	info, err := botAPI.GetWebhookInfo()
	if err != nil {
		return err
	}
	if got, want := info.URL, ""; got != want {
		return errors.New("webhook was not destroyed")
	}
	return nil
}
