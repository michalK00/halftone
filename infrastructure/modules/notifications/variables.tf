variable "environment" {
  description = "Environment name"
  type        = string
}

variable "project_name" {
  description = "Name of the project"
  type        = string
}

# SNS Configuration
variable "sns_topic_name" {
  description = "Name of the SNS topic for email notifications"
  type        = string
}

variable "email_subscriber" {
  description = "Email address to subscribe to notifications"
  type        = string
}

variable "sns_display_name" {
  description = "Display name for SNS topic"
  type        = string
}

# KMS Configuration
variable "enable_kms_encryption" {
  description = "Enable KMS encryption for SNS"
  type        = bool
  default     = true
}

variable "sns_kms_master_key_id" {
  description = "KMS key ID for SNS encryption (if not provided, AWS managed key will be used)"
  type        = string
  default     = "alias/aws/sns"
}