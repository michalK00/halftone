resource "aws_iam_user" "developer" {
  name = "${var.project_name}-${var.environment}-developer"

  tags = {
    Name        = "${var.project_name}-${var.environment}-developer"
    Environment = var.environment
  }
}

resource "aws_iam_access_key" "developer" {
  user = aws_iam_user.developer.name
}

resource "aws_iam_policy" "s3_access" {
  name        = "${var.project_name}-${var.environment}-s3-access"
  description = "Policy for accessing S3 buckets"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ]
        Resource = [
          var.photos_bucket_arn,
          "${var.photos_bucket_arn}/*",
          var.logs_bucket_arn,
          "${var.logs_bucket_arn}/*"
        ]
      }
    ]
  })
}

resource "aws_iam_policy" "cognito_access" {
  name        = "${var.project_name}-${var.environment}-cognito-access"
  description = "Policy for accessing Cognito"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "cognito-idp:AdminCreateUser",
          "cognito-idp:AdminDeleteUser",
          "cognito-idp:AdminUpdateUserAttributes",
          "cognito-idp:AdminGetUser",
          "cognito-idp:ListUsers"
        ]
        Resource = var.user_pool_arn
      }
    ]
  })
}

resource "aws_iam_user_policy_attachment" "developer_s3" {
  user       = aws_iam_user.developer.name
  policy_arn = aws_iam_policy.s3_access.arn
}

resource "aws_iam_user_policy_attachment" "developer_cognito" {
  user       = aws_iam_user.developer.name
  policy_arn = aws_iam_policy.cognito_access.arn
}

# Future ECS task execution role (prepare for later)
resource "aws_iam_role" "ecs_task_execution" {
  name = "${var.project_name}-${var.environment}-ecs-task-execution"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ecs-tasks.amazonaws.com"
      }
    }]
  })

  tags = {
    Name        = "${var.project_name}-${var.environment}-ecs-task-execution"
    Environment = var.environment
  }
}

# Attach AWS managed policy for ECS task execution
resource "aws_iam_role_policy_attachment" "ecs_task_execution" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# Future ECS task role (for application permissions)
resource "aws_iam_role" "ecs_task" {
  name = "${var.project_name}-${var.environment}-ecs-task"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ecs-tasks.amazonaws.com"
      }
    }]
  })

  tags = {
    Name        = "${var.project_name}-${var.environment}-ecs-task"
    Environment = var.environment
  }
}

# Attach S3 and Cognito policies to ECS task role
resource "aws_iam_role_policy_attachment" "ecs_task_s3" {
  role       = aws_iam_role.ecs_task.name
  policy_arn = aws_iam_policy.s3_access.arn
}

resource "aws_iam_role_policy_attachment" "ecs_task_cognito" {
  role       = aws_iam_role.ecs_task.name
  policy_arn = aws_iam_policy.cognito_access.arn
}