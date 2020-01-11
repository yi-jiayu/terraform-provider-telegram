# Telegram Provider

The Telegram provider is used to interact with Telegram-related resources.

## Example Usage

```hcl
variable "bot_token" {
  type = string
}

provider "telegram" {
  bot_token = var.bot_token
}
```

## Argument Reference

The following arguments are supported in the Telegram `provider` block:

* `bot_token` (Required) - The unique authentication token provided when a bot is created.
