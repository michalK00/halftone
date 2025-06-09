output "cluster_id" {
  description = "ECS cluster ID"
  value       = aws_ecs_cluster.main.id
}

output "cluster_name" {
  description = "ECS cluster name"
  value       = aws_ecs_cluster.main.name
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