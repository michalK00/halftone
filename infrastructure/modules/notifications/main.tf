locals {
  common_tags = merge({
    Environment = var.environment
    Project     = var.project_name
    Module      = "notifications"
  })
}

resource "aws_sns_topic" "email_notifications" {
  name         = var.sns_topic_name
  display_name = var.sns_display_name != "" ? var.sns_display_name : "Email Notifications - ${var.project_name}"

  kms_master_key_id = var.enable_kms_encryption ? var.sns_kms_master_key_id : null

  delivery_policy = jsonencode({
    "http" : {
      "defaultHealthyRetryPolicy" : {
        "minDelayTarget" : 20,
        "maxDelayTarget" : 20,
        "numRetries" : 3,
        "numMaxDelayRetries" : 0,
        "numMinDelayRetries" : 0,
        "numNoDelayRetries" : 0,
        "backoffFunction" : "linear"
      },
      "disableSubscriptionOverrides" : false,
      "defaultThrottlePolicy" : {
        "maxReceivesPerSecond" : 1
      }
    }
  })

  tags = merge(local.common_tags, {
    Name = var.sns_topic_name
    Type = "email-notifications"
  })
}

# SNS Topic Policy
resource "aws_sns_topic_policy" "email_notifications" {
  arn    = aws_sns_topic.email_notifications.arn
  policy = data.aws_iam_policy_document.sns_topic_policy.json
}

resource "aws_sns_topic_subscription" "email_notifications" {

  topic_arn = aws_sns_topic.email_notifications.arn
  protocol  = "email"
  endpoint  = var.email_subscriber

  # Email subscriptions need to be confirmed manually
  confirmation_timeout_in_minutes = 5
}

