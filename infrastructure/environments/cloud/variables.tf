variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "eu-north-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "prod"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "halftone"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "List of availability zones"
  type        = list(string)
  default     = ["eu-north-1a", "eu-north-1b"]
}

variable "api_image_tag" {
  description = "API image tag"
  type        = string
  default     = "latest"
}

variable "client_image_tag" {
  description = "Client image tag"
  type        = string
  default     = "latest"
}

variable "admin_image_tag" {
  description = "Admin image tag"
  type        = string
  default     = "latest"
}

variable "api_desired_count" {
  description = "API task count"
  type        = number
  default     = 1
}

variable "frontend_desired_count" {
  description = "Frontend task count"
  type        = number
  default     = 1
}

variable "docdb_master_username" {
  description = "DocumentDB username"
  type        = string
  sensitive   = true
}

variable "docdb_master_password" {
  description = "DocumentDB password"
  type        = string
  sensitive   = true
}

variable "mongodb_database_name" {
  description = "Database name"
  type        = string
  default     = "halftone"
}

variable "notification_email" {
  description = "Email for CloudWatch alarms"
  type        = string
  default     = "mklemens49@gmail.com"
}

variable "lambda_source_dir" {
  description = "Directory containing the Lambda source code"
  type        = string
  default     = "../../../application/lambda/image-processing"
}

variable "google_application_credentials" {
  description = "Google application credentials for FCM"
  type        = string
  sensitive = true
}

variable "fcm_project_id" {
  description = "FCM project ID for push notifications"
  type        = string
  sensitive = true
}