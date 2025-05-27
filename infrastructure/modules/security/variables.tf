variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "photos_bucket_arn" {
  description = "ARN of the photos S3 bucket"
  type        = string
}

variable "logs_bucket_arn" {
  description = "ARN of the logs S3 bucket"
  type        = string
}

variable "user_pool_arn" {
  description = "ARN of the Cognito User Pool"
  type        = string
}