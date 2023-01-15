package telegram

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/yi-jiayu/ted"
)

func TestAccResourceTelegramBotCommands(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceTelegramBotCommandsDestroy,
		Steps: []resource.TestStep{
			{
				Config: `resource "telegram_bot_commands" "example" {
  commands = [
    {
      command = "start",
      description = "View welcome message"
    },
  ]
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("telegram_bot_commands.example", "commands.#", "1"),
					resource.TestCheckResourceAttr("telegram_bot_commands.example", "commands.0.command", "start"),
					resource.TestCheckResourceAttr("telegram_bot_commands.example", "commands.0.description", "View welcome message"),
					testAccResourceTelegramBotCommands("telegram_bot_commands.example"),
				),
			},
			{
				Config: `resource "telegram_bot_commands" "example" {
  commands = [
    {
      command = "start",
      description = "View welcome message"
    },
    {
      command = "help",
      description = "Show help"
    }
  ]
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("telegram_bot_commands.example", "commands.#", "2"),
					resource.TestCheckResourceAttr("telegram_bot_commands.example", "commands.0.command", "start"),
					resource.TestCheckResourceAttr("telegram_bot_commands.example", "commands.0.description", "View welcome message"),
					resource.TestCheckResourceAttr("telegram_bot_commands.example", "commands.1.command", "help"),
					resource.TestCheckResourceAttr("telegram_bot_commands.example", "commands.1.description", "Show help"),
					testAccResourceTelegramBotCommands("telegram_bot_commands.example"),
				),
			},
		},
	})
}

func testAccResourceTelegramBotCommands(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}
		botAPI := testAccProvider.Meta().(ted.Bot)
		commands, err := botAPI.GetMyCommands()
		if err != nil {
			return err
		}
		attrs := rs.Primary.Attributes
		path := "commands.#"
		if got, want := attrs[path], strconv.Itoa(len(commands)); got != want {
			return fmt.Errorf("wanted %s to be %s, got %s", path, want, got)
		}
		for i, command := range commands {
			var path string
			path = fmt.Sprintf("commands.%d.%s", i, "command")
			if got, want := attrs[path], command.Command; got != want {
				return fmt.Errorf("wanted %s to be %s, got %s", path, want, got)
			}
			path = fmt.Sprintf("commands.%d.%s", i, "description")
			if got, want := attrs[path], command.Description; got != want {
				return fmt.Errorf("wanted %s to be %s, got %s", path, want, got)
			}
		}
		return nil
	}
}

func testAccResourceTelegramBotCommandsDestroy(*terraform.State) error {
	botAPI := testAccProvider.Meta().(ted.Bot)
	commands, err := botAPI.GetMyCommands()
	if err != nil {
		return err
	}
	if len(commands) != 0 {
		return errors.New("commands was not destroyed")
	}
	return nil
}
