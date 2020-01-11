variable "bot_token" {
  type = string
}

provider "telegram" {
  bot_token = var.bot_token
}

resource "telegram_bot_webhook" "example" {
  url = "https://example.com/webhook"
}
