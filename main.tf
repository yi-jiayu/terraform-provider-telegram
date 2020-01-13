variable "bot_token" {
  type = string
}

provider "telegram" {
  bot_token = var.bot_token
}

data "telegram_bot" "example" {}

output "bot_link" {
  value = "https://t.me/${data.telegram_bot.example.username}"
}

resource "telegram_bot_webhook" "example" {
  url = "https://example.com/webhook"
  certificate = file("cert.pem")
}
