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

output "ecs_task_execution_role_arn" {
  description = "ARN of the ECS task execution role"
  value       = module.security.ecs_task_execution_role_arn
}

output "ecs_task_role_arn" {
  description = "ARN of the ECS task role"
  value       = module.security.ecs_task_role_arn
}

output "database_security_group_id" {
  description = "ID of the database security group"
  value       = module.networking.database_security_group_id
}

output "db_subnet_group_name" {
  description = "Name of the DB subnet group"
  value       = module.networking.db_subnet_group_name
}

output "alb_dns_name" {
  description = "DNS name of the load balancer"
  value       = module.ecs.alb_dns_name
}

output "application_urls" {
  description = "Application URLs"
  value = {
    client = "https://${module.ecs.alb_dns_name}"
    api    = "https://${module.ecs.alb_dns_name}/api"
  }
}

output "ecs_cluster_name" {
  description = "ECS cluster name"
  value       = module.ecs.ecs_cluster_name
}

output "codedeploy_app_name" {
  description = "CodeDeploy application name"
  value       = module.ecs.codedeploy_app_name
}

output "codedeploy_deployment_group_name" {
  description = "CodeDeploy deployment group name"
  value = module.ecs.codedeploy_deployment_group_name
}

output "task_definition_family" {
  description = "ECS task definition family"
  value = module.ecs.task_definition_family
}

