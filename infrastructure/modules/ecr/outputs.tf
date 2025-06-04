output "api_repository_url" {
  description = "URL of the API (backend) ECR repository"
  value       = aws_ecr_repository.api.repository_url
}

output "client_repository_url" {
  description = "URL of the client frontend ECR repository"
  value       = aws_ecr_repository.client.repository_url
}

output "registry_id" {
  description = "The registry ID where the repositories were created"
  value       = aws_ecr_repository.api.registry_id
}