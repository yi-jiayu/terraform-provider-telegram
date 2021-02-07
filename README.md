# Terraform Provider for Telegram
Manage Telegram resources with Terraform

## Example Usage

If we have a bot implemented as a Google Cloud Function, we can set the bot
webhook in the same Terraform project:

```hcl
resource "google_cloudfunctions_function" "my_bot" {
  name         = "my-bot"
  trigger_http = true

  # other arguments
}

resource "telegram_bot_webhook" "my_bot" {
  url = google_cloudfunctions_function.my_bot.https_trigger_url
}
```

You can also manage bot commands using Terraform:

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

## Installation

You can find this provider on the [Terraform
Registry](https://registry.terraform.io/):

https://registry.terraform.io/providers/yi-jiayu/telegram/latest

## Reference

### [Telegram provider](website/docs/index.html.markdown)

### Data Sources

- [`telegram_bot`](website/docs/d/bot.html.markdown) - Get information about
  the currently-authenticated Telegram bot

### Resources

- [`telegram_bot_webhook`](website/docs/r/bot_webhook.html.markdown) - Manage
  the webhook for the currently-authenticated Telegram bot
- [`telegram_bot_commands`](website/docs/r/bot_commands.html.markdown) - Manage
  commands for the currently-authenticated Telegram bot
