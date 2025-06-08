variable "environment" {
    description = "Environment name"
    type        = string
}

variable "project_name" {
    description = "Project name"
    type        = string
}

variable "vpc_id" {
    description = "VPC ID"
    type        = string
}

variable "subnet_ids" {
    description = "List of subnet IDs"
    type        = list(string)
}

variable "app_security_group_ids" {
    description = "List of security group IDs for the application"
    type        = list(string)
}

variable "master_username" {
    description = "Master username for the database"
    type        = string
}

variable "master_password" {
    description = "Master password for the database"
    type        = string
    sensitive   = true
}

variable "instance_count" {
    description = "Number of database instances"
    type        = number
    default     = 1
}

variable "instance_type" {
  type = string
  description = "Instance type for the database"
}
