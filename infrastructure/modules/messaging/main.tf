locals {
  dlq_name       = "${var.sqs_queue_name}-dlq"

  common_tags = merge({
    Environment = var.environment
    Project     = var.project_name
    Module      = "notifications"
  })
}

# Dead Letter Queue (created first as it's referenced by main queue)
resource "aws_sqs_queue" "lambda_triggers_dlq" {

  name = local.dlq_name

  message_retention_seconds = var.dlq_message_retention_seconds

  kms_master_key_id                 = var.enable_kms_encryption ? var.kms_master_key_id : null
  kms_data_key_reuse_period_seconds = var.enable_kms_encryption ? 300 : null

  tags = merge(local.common_tags, {
    Name = local.dlq_name
    Type = "dead-letter-queue"
  })
}

resource "aws_sqs_queue" "lambda_triggers" {

  name = var.sqs_queue_name

  delay_seconds             = var.delay_seconds
  max_message_size         = 262144 # 256 KB
  message_retention_seconds = var.message_retention_seconds
  receive_wait_time_seconds = var.receive_wait_time_seconds
  visibility_timeout_seconds = var.visibility_timeout_seconds

  # Dead Letter Queue Configuration
  redrive_policy = var.enable_dlq ? jsonencode({
    deadLetterTargetArn = aws_sqs_queue.lambda_triggers_dlq.arn
    maxReceiveCount     = var.max_receive_count
  }) : null

  kms_master_key_id                 = var.enable_kms_encryption ? var.kms_master_key_id : null
  kms_data_key_reuse_period_seconds = var.enable_kms_encryption ? 300 : null

  tags = merge(local.common_tags, {
    Name = var.sqs_queue_name
    Type = "lambda-triggers"
  })
}

resource "aws_sqs_queue_policy" "lambda_triggers" {
  queue_url = aws_sqs_queue.lambda_triggers.id
  policy    = data.aws_iam_policy_document.sqs_queue_policy.json
}

# CloudWatch Alarms for Monitoring
resource "aws_cloudwatch_metric_alarm" "sqs_dlq_messages" {
  alarm_name          = "${local.dlq_name}-messages-alarm"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "ApproximateNumberOfVisibleMessages"
  namespace           = "AWS/SQS"
  period              = "300"
  statistic           = "Average"
  threshold           = "0"
  alarm_description   = "This metric monitors messages in the dead letter queue"
  alarm_actions       = [var.aws_sns_topic_arn]

  dimensions = {
    QueueName = local.dlq_name
  }

  tags = local.common_tags
}

resource "aws_cloudwatch_metric_alarm" "sqs_age_of_oldest_message" {
  alarm_name          = "${local.dlq_name}-age-alarm"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "ApproximateAgeOfOldestMessage"
  namespace           = "AWS/SQS"
  period              = "300"
  statistic           = "Maximum"
  threshold           = "300" # 5 minutes
  alarm_description   = "This metric monitors the age of the oldest message in the queue"
  alarm_actions       = [var.aws_sns_topic_arn]


  dimensions = {
    QueueName = var.sqs_queue_name
  }

  tags = local.common_tags
}