# Data Source: telegram_bot

Use this data source to get details about the currently authenticated Telegram
bot.

## Example Usage

```hcl

data "telegram_bot" "example" {}

output "bot_link" {
  value = "https://t.me/${data.telegram_bot.example.username}"
}
```


## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

* `id` - Unique identifier for the bot.
* `name` - Bot name.
* `username` - Bot username.
