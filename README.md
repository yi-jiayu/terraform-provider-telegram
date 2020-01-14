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

## Installation

Download the [latest
release](https://github.com/yi-jiayu/terraform-provider-telegram/releases) for
your OS and extract it to the user plugins directory, which is located at
`%APPDATA%\terraform.d\plugins` on Windows and `~/.terraform.d/plugins` on
Linux and macOS.

See also the [Terraform documentation on third-party
providers](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

## Reference

### [Telegram provider](website/docs/index.html.markdown)

### Data Sources

- [`telegram_bot`](website/docs/d/bot.html.markdown) - Get information about
  the currently-authenticated Telegram bot

### Resources

- [`telegram_bot_webhook`](website/docs/r/bot_webhook.html.markdown) - Manage
  the webhook for the currently-authenticated Telegram bot
