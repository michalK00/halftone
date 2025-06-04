output "vpc_id" {
  description = "ID of the VPC"
  value       = module.networking.vpc_id
}

output "public_subnet_ids" {
  description = "IDs of public subnets"
  value       = module.networking.public_subnet_ids
}

output "photos_bucket_name" {
  description = "Name of the photos S3 bucket"
  value       = module.storage.photos_bucket_name
}

output "user_pool_id" {
  description = "ID of the Cognito User Pool"
  value       = module.auth.user_pool_id
}

output "user_pool_client_id" {
  description = "ID of the Cognito User Pool Client"
  value       = module.auth.user_pool_client_id
}

# IAM Outputs (sensitive)
output "developer_access_key_id" {
  description = "Access key ID for developer user"
  value       = module.security.developer_access_key_id
  sensitive   = true
}

output "developer_secret_access_key" {
  description = "Secret access key for developer user"
  value       = module.security.developer_secret_access_key
  sensitive   = true
}

# Future ECS Role ARNs
output "ecs_task_execution_role_arn" {
  description = "ARN of the ECS task execution role"
  value       = module.security.ecs_task_execution_role_arn
}

output "ecs_task_role_arn" {
  description = "ARN of the ECS task role"
  value       = module.security.ecs_task_role_arn
}

# Database preparation outputs
output "database_security_group_id" {
  description = "ID of the database security group"
  value       = module.networking.database_security_group_id
}

output "db_subnet_group_name" {
  description = "Name of the DB subnet group"
  value       = module.networking.db_subnet_group_name
}