terraform {
  required_version = ">= 1.4"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }

  backend "s3" {
    encrypt = true
    region  = "eu-central-1"
    key     = "telegram-bot"
  }
}

# Configure the AWS Provider
provider "aws" {
  region = "eu-central-1"
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

resource "aws_lambda_function" "lambda_handler" {

  function_name = "${var.app_name}-${var.stack_env}"
  description   = "Telegram bot with AWS Lambda"

  publish                        = true
  package_type                   = "Image"
  image_uri                      = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${data.aws_region.current.name}.amazonaws.com/${var.app_name}:${var.app_version}"
  role                           = aws_iam_role.lambda.arn
  memory_size                    = var.memory
  timeout                        = var.timeout
  reserved_concurrent_executions = var.reserved_concurrency

  image_config {
    command = ["/app/service"]
  }

  environment {
    variables = {
      S3_BUCKET_REGION = "eu-central-1"
    }
  }

  tags = {
    Name = "${var.app_name}-${var.stack_env}"
  }
}

resource "aws_cloudwatch_log_group" "log_group" {
  name = "/aws/lambda/${aws_lambda_function.lambda_handler.function_name}"

  retention_in_days = 7
  tags = {
    Environment = var.stack_env
    Service     = var.app_name
  }
}

data "aws_iam_policy_document" "lambda" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "lambda" {
  name               = "${var.app_name}-${var.stack_env}"
  assume_role_policy = data.aws_iam_policy_document.lambda.json
}

resource "aws_iam_role_policy_attachment" "execution_role" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_function_url" "url" {
  function_name      = aws_lambda_function.lambda_handler.function_name
  authorization_type = "NONE"
}

resource "aws_lambda_permission" "lambda_invoke_url" {
  statement_id           = "FunctionURLAllowPublicAccess"
  action                 = "lambda:InvokeFunctionUrl"
  function_name          = aws_lambda_function.lambda_handler.function_name
  principal              = "*"
  function_url_auth_type = "NONE"

}