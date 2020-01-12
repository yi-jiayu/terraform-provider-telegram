package telegram

import (
	"fmt"
	"strconv"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceTelegramBot(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "telegram_bot" "example" {}`,
				Check:  testAccDataSourceTelegramBot("data.telegram_bot.example"),
			},
		},
	})
}

func testAccDataSourceTelegramBot(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}
		bot := testAccProvider.Meta().(*tgbotapi.BotAPI).Self
		if got, want := rs.Primary.ID, strconv.Itoa(bot.ID); got != want {
			return fmt.Errorf("wanted ID to be %s, got %s", want, got)
		}
		if got, want := rs.Primary.Attributes["user_id"], strconv.Itoa(bot.ID); got != want {
			return fmt.Errorf("wanted user_id to be %s, got %s", want, got)
		}
		if got, want := rs.Primary.Attributes["name"], bot.FirstName; got != want {
			return fmt.Errorf("wanted name to be %s, got %s", want, got)
		}
		if got, want := rs.Primary.Attributes["username"], bot.UserName; got != want {
			return fmt.Errorf("wanted username to be %s, got %s", want, got)
		}
		return nil
	}
}
