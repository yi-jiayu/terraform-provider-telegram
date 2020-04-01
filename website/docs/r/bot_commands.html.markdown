# Resource: telegram_bot_commands

Set commands for a Telegram Bot

For more information, refer to:

* [API reference](https://core.telegram.org/bots/api#setMyCommands)

## Example Usage

```hcl
resource "telegram_bot_commands" "example" {
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
}
```

## Argument Reference

The following arguments are supported:

* `commands` - (Required) The list of commands for the bot. Structure is documented below.

The `commands` block supports:

* `command` (Required) - The name of the command
* `description` (Required) - Help text for the command 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:
