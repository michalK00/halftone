variable "environment" {
  description = "Environment name"
  type        = string
}

variable "project_name" {
  description = "Name of the project"
  type        = string
}

variable "aws_sns_topic_arn" {
    description = "ARN of the SNS topic for email notifications"
    type        = string
}

variable "sqs_queue_name" {
  description = "Name of the SQS queue for Lambda triggers"
  type        = string
}

variable "message_retention_seconds" {
  description = "The number of seconds Amazon SQS retains a message"
  type        = number
  default     = 3600 # 1 hour
}

variable "visibility_timeout_seconds" {
  description = "The visibility timeout for the queue"
  type        = number
  default     = 120
}

variable "max_receive_count" {
  description = "Maximum number of times a message can be received before being moved to DLQ"
  type        = number
  default     = 3
}

variable "delay_seconds" {
  description = "The time in seconds that the delivery of all messages in the queue will be delayed"
  type        = number
  default     = 0
}

variable "receive_wait_time_seconds" {
  description = "The time for which a ReceiveMessage call will wait for a message to arrive"
  type        = number
  default     = 0
}

variable "enable_dlq" {
  description = "Enable Dead Letter Queue"
  type        = bool
  default     = true
}

variable "dlq_message_retention_seconds" {
  description = "Message retention period for Dead Letter Queue"
  type        = number
  default     = 604800 # 7 days
}

variable "enable_kms_encryption" {
  description = "Enable KMS encryption for SQS and SNS"
  type        = bool
  default     = true
}

variable "kms_master_key_id" {
  description = "KMS key ID for encryption (if not provided, AWS managed key will be used)"
  type        = string
  default     = "alias/aws/sqs"
}
