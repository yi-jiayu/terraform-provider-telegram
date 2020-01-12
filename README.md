# Terraform Provider for Telegram

Currently only supports bot-related resources. Useful in Telegram bot projects.

## Example Usage

If we have a bot implemented as a Google Cloud Function, we can set the bot webhook in the same Terraform project:

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
