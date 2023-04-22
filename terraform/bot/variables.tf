variable "stack_env" {
  type        = string
  description = "Stack environment"
  default     = "prod"
}

variable "app_name" {
  type        = string
  description = "Application name"
  default     = "telegram-bot"
}

variable "app_version" {
  type        = string
  default     = "v0.0.5"
  description = "Container image version used to deploy the lambda function"
}

variable "memory" {
  type        = string
  default     = "128"
  description = "Memory in MB to assign to the Lambda."
}

variable "timeout" {
  type        = string
  default     = "900"
  description = "Seconds in which the Lambda should run before timing out."
}

variable "reserved_concurrency" {
  type        = number
  description = "Reserved concurrency guarantees the maximum number of concurrent instances for the function. A value of 0 disables lambda from being triggered and -1 removes any concurrency limitations. "
  default     = -1
}

variable "telegram_bot_token" {
  type        = string
  description = "Telegram bot token"
  default     = ""
}

variable "telegram_bot_webhook" {
    type        = string
    description = "Telegram bot webhook"
    default     = ""
}