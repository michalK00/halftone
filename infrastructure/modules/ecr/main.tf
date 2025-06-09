resource "aws_ecr_repository" "api" {
  name                 = "halftone-api"
  image_tag_mutability = "MUTABLE"

  tags = {
    Name        = "halftone-api"
    Environment = var.environment
    Service     = "api"
  }
}

resource "aws_ecr_repository" "client" {
  name                 = "halftone-client"
  image_tag_mutability = "MUTABLE"

  tags = {
    Name        = "halftone-client"
    Environment = var.environment
    Service     = "client"
  }
}

locals {
  repositories = {
    api    = aws_ecr_repository.api
    client = aws_ecr_repository.client
  }
}

resource "aws_ecr_lifecycle_policy" "main" {
  for_each   = local.repositories
  repository = each.value.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep last 3 images"
        selection = {
          tagStatus     = "tagged"
          tagPrefixList = ["v"]
          countType     = "imageCountMoreThan"
          countNumber   = 3
        }
        action = {
          type = "expire"
        }
      },
      {
        rulePriority = 2
        description  = "Remove untagged images after 2 days"
        selection = {
          tagStatus   = "untagged"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 2
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}
