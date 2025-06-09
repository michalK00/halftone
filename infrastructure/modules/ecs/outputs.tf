output "cluster_id" {
  description = "ECS cluster ID"
  value       = aws_ecs_cluster.main.id
}

output "ecs_cluster_name" {
  description = "ECS cluster name"
  value       = aws_ecs_cluster.main.name
}

output "codedeploy_app_name" {
  description = "CodeDeploy application name"
  value       = aws_codedeploy_app.api.name
}

output "codedeploy_deployment_group_name" {
  description = "CodeDeploy deployment group name"
  value = aws_codedeploy_deployment_group.api.deployment_group_name
}

output "task_definition_family" {
  description = "ECS task definition family"
  value = aws_ecs_task_definition.api.family
}


output "alb_dns_name" {
  description = "ALB DNS name"
  value       = aws_lb.main.dns_name
}

output "target_group_blue_name" {
  description = "Blue target group name"
  value       = aws_lb_target_group.api.name
}

output "target_group_green_name" {
  description = "Green target group name"
  value       = aws_lb_target_group.api_green.name
}