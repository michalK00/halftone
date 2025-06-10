variable "environment" {
  description = "Environment name"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "public_subnet_ids" {
  description = "Public subnet IDs"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "Private subnet IDs"
  type        = list(string)
}

variable "alb_security_group_id" {
  description = "ALB security group ID"
  type        = string
}

variable "ecs_tasks_security_group_id" {
  description = "ECS security group ID"
  type        = string
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.small"
}

variable "api_image" {
  description = "API container image URL"
  type        = string
}

variable "api_image_tag" {
  description = "API container image tag"
  type        = string
  default     = "latest"
}

variable "client_image" {
  description = "Client container image URL"
  type        = string
}

variable "client_image_tag" {
  description = "Client container image tag"
  type        = string
  default     = "latest"
}

variable "api_desired_count" {
  description = "Desired number of API tasks"
  type        = number
  default     = 1
}

variable "frontend_desired_count" {
  description = "Desired number of frontend tasks"
  type        = number
  default     = 1
}

variable "cognito_user_pool_id" {
  description = "Cognito User Pool ID"
  type        = string
}

variable "cognito_app_client_id" {
  description = "Cognito App Client ID"
  type        = string
}

variable "cognito_app_client_secret_arn" {
  description = "Cognito App Client Secret ARN"
  type        = string
  default     = ""
}

variable "s3_name" {
  description = "S3 bucket name"
  type        = string
}

variable "s3_uri" {
  description = "S3 bucket URI"
  type        = string
}

variable "mongodb_uri_arn" {
  description = "MongoDB URI secret ARN"
  type        = string
  default     = ""
}

variable "mongodb_database_name" {
  description = "MongoDB database name"
  type        = string
}

variable "sqs_queue_name" {
  description = "SQS queue name for Lambda triggers"
  type        = string
}

variable "sqs_queue_url" {
  description = "SQS queue url for Lambda triggers"
  type        = string
}

variable "sqs_queue_arn" {
  description = "ARN of the SQS queue for Lambda triggers"
  type        = string
}

variable "desired_capacity" {
  description = "Desired capacity for ECS cluster"
  type        = number
  default = 2
}

variable "max_size" {
  description = "Maximum size for ECS cluster"
  type        = number
  default = 4
}

variable "min_size" {
  description = "Minimum size for ECS cluster"
  type        = number
  default = 2
}

variable "fcm_project_id" {
  description = "FCM project ID for push notifications"
  type        = string
  sensitive = true
}

variable "google_application_credentials" {
  description = "Google application credentials for FCM"
  type        = string
  sensitive = true
}