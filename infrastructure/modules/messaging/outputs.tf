output "sqs_queue_arn" {
  description = "ARN of the SQS queue for Lambda triggers"
  value       = aws_sqs_queue.lambda_triggers.arn
}

output "sqs_queue_url" {
  description = "URL of the SQS queue"
  value       = aws_sqs_queue.lambda_triggers.url
}

output "sqs_queue_name" {
  description = "Name of the SQS queue"
  value       = aws_sqs_queue.lambda_triggers.name
}

output "sqs_queue_id" {
  description = "ID of the SQS queue"
  value       = aws_sqs_queue.lambda_triggers.id
}

output "sqs_dlq_arn" {
  description = "ARN of the SQS dead letter queue"
  value       = var.enable_dlq ? aws_sqs_queue.lambda_triggers_dlq.arn : null
}

output "sqs_dlq_url" {
  description = "URL of the SQS dead letter queue"
  value       = var.enable_dlq ? aws_sqs_queue.lambda_triggers_dlq.url : null
}

output "sqs_dlq_name" {
  description = "Name of the SQS dead letter queue"
  value       = var.enable_dlq ? aws_sqs_queue.lambda_triggers_dlq.name : null
}

output "cloudwatch_alarm_arns" {
  description = "ARNs of CloudWatch alarms created for monitoring"
  value = {
    dlq_messages_alarm     = var.enable_dlq ? aws_cloudwatch_metric_alarm.sqs_dlq_messages.arn : null
    age_of_oldest_message  = aws_cloudwatch_metric_alarm.sqs_age_of_oldest_message.arn
  }
}

output "messaging_resources" {
  description = "Complete information about messaging resources"
  value = {

    sqs = {
      queue = {
        arn  = aws_sqs_queue.lambda_triggers.arn
        url  = aws_sqs_queue.lambda_triggers.url
        name = aws_sqs_queue.lambda_triggers.name
      }
      dlq = var.enable_dlq ? {
        arn  = aws_sqs_queue.lambda_triggers_dlq.arn
        url  = aws_sqs_queue.lambda_triggers_dlq.url
        name = aws_sqs_queue.lambda_triggers_dlq.name
      } : null
    }
  }
}