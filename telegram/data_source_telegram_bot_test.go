package telegram

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/yi-jiayu/ted"
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
		bot, err := testAccProvider.Meta().(ted.Bot).GetMe()
		if err != nil {
			return err
		}
		if got, want := rs.Primary.ID, strconv.FormatInt(bot.ID, 10); got != want {
			return fmt.Errorf("wanted ID to be %s, got %s", want, got)
		}
		if got, want := rs.Primary.Attributes["user_id"], strconv.FormatInt(bot.ID, 10); got != want {
			return fmt.Errorf("wanted user_id to be %s, got %s", want, got)
		}
		if got, want := rs.Primary.Attributes["name"], bot.FirstName; got != want {
			return fmt.Errorf("wanted name to be %s, got %s", want, got)
		}
		if got, want := rs.Primary.Attributes["username"], bot.Username; got != want {
			return fmt.Errorf("wanted username to be %s, got %s", want, got)
		}
		return nil
	}
}
