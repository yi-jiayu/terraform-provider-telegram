# Resource: telegram_bot_webhook

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
* `max_connections` - (Optional) Maximum allowed number of simultaneous
  connections to the webhook for update delivery, 1-100. Defaults to 40. Use
  lower values to limit the load on your bot‘s server, and higher values to
  increase your bot’s throughput.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:
