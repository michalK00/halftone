variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "project_name" {
  description = "Project name"
  type        = string
}

variable "function_name" {
  description = "Function name"
  type        = string
  default     = "photo-processor"
}

variable "sqs_queue_arn" {
  description = "ARN of the SQS queue to trigger the Lambda"
  type        = string
}

variable "sns_topic_arn" {
  description = "SNS topic ARN for alarms"
  type        = string
}

variable "tags" {
  description = "Tags for resources"
  type        = map(string)
  default     = {}
}

variable "photos_bucket_arn" {
    description = "ARN of the S3 bucket for photos"
    type        = string
}

# Source code for the Lambda function
variable "source_dir" {
  description = "Directory containing the Lambda source code"
  type        = string
}

variable "binary_name" {
  description = "Name of the compiled binary"
  type        = string
  default     = "bootstrap"
}

variable "build_command" {
  description = "Command to build the Lambda function"
  type        = string
  default     = "GOARCH=amd64 GOOS=linux go build -o bootstrap main.go"
}