locals {
  common_tags = merge(var.tags, {
    Environment = var.environment
    Project     = var.project_name
  })
}

resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${var.function_name}"
  retention_in_days = 14

  tags = local.common_tags
}

resource "aws_iam_role" "lambda_role" {
  name = "${var.function_name}-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = local.common_tags
}

resource "aws_iam_role_policy" "lambda_policy" {
  name = "${var.function_name}-policy"
  role = aws_iam_role.lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "arn:aws:logs:*:*:*"
      },
      {
        Effect = "Allow"
        Action = [
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes"
        ]
        Resource = var.sqs_queue_arn
      }
    ]
  })
}

resource "aws_iam_role_policy" "lambda_s3_policy" {
  name = "${var.project_name}-${var.environment}-lambda-s3"
  role = aws_iam_role.lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ]
        Resource = "${var.photos_bucket_arn}/*"
      },
      {
        Effect = "Allow"
        Action = [
          "s3:ListBucket"
        ]
        Resource = var.photos_bucket_arn
      }
    ]
  })
}

resource "null_resource" "lambda_build" {
  # triggers = {
  #   source_hash = filebase64sha256("${var.source_dir}/main.go")
  # }

  provisioner "local-exec" {
    command     = var.build_command
    working_dir = var.source_dir
  }
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "${var.source_dir}/${var.binary_name}"
  output_path = "${var.source_dir}/lambda-deployment.zip"

  depends_on = [null_resource.lambda_build]
}

resource "aws_lambda_function" "photo_processor" {
  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  function_name    = var.function_name
  role            = aws_iam_role.lambda_role.arn
  handler         = "main.handler"
  runtime         = "provided.al2023"
  timeout         = 60 * 2
  memory_size     = 128

  tags = local.common_tags

  depends_on = [aws_cloudwatch_log_group.lambda_logs]
}

resource "aws_lambda_event_source_mapping" "sqs_trigger" {
  event_source_arn = var.sqs_queue_arn
  function_name    = aws_lambda_function.photo_processor.arn
  batch_size       = 10
  enabled          = true
}

resource "aws_cloudwatch_metric_alarm" "lambda_errors" {
  alarm_name          = "${var.function_name}-errors"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = "300"
  statistic           = "Sum"
  threshold           = "0"
  alarm_description   = "Lambda function errors"

  dimensions = {
    FunctionName = aws_lambda_function.photo_processor.function_name
  }

  alarm_actions = [var.sns_topic_arn]

  tags = local.common_tags
}

resource "aws_cloudwatch_metric_alarm" "lambda_duration" {
  alarm_name          = "${var.function_name}-duration"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "Duration"
  namespace           = "AWS/Lambda"
  period              = "300"
  statistic           = "Average"
  threshold           = "120000"  # 120 seconds
  alarm_description   = "Lambda function duration"

  dimensions = {
    FunctionName = aws_lambda_function.photo_processor.function_name
  }

  alarm_actions = [var.sns_topic_arn]

  tags = local.common_tags
}