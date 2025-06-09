output "sns_topic_arn" {
  description = "ARN of the SNS topic for email notifications"
  value       = aws_sns_topic.email_notifications.arn
}

output "sns_topic_name" {
  description = "Name of the SNS topic"
  value       = aws_sns_topic.email_notifications.name
}

output "sns_topic_id" {
  description = "ID of the SNS topic"
  value       = aws_sns_topic.email_notifications.id
}

output "notification_resources" {
  description = "Complete information about notification resources"
  value = {
    sns = {
      arn           = aws_sns_topic.email_notifications.arn
      name          = aws_sns_topic.email_notifications.name
      display_name  = aws_sns_topic.email_notifications.display_name
      subscriber   = var.email_subscriber
    }
  }
}