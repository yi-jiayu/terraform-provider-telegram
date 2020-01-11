# Resource: bot_webhook

Sets up a webhook for a Telegram bot.

For more information, refer to:

* [Official documentation](https://core.telegram.org/bots/webhooks)

* [API reference](https://core.telegram.org/bots/api#setwebhook)

## Example Usage

```hcl
resource "telegram_bot_webhook" "example" {
  url = "https://example.com/webhook"
}
```

## Argument Reference

The following arguments are supported:

* `url` - (Required) The webhook URL.
