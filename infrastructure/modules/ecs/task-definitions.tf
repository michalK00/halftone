resource "aws_ecs_task_definition" "api" {
  family                   = "${var.environment}-api"
  network_mode             = "bridge"
  requires_compatibilities = ["EC2"]
  cpu                      = "512"
  memory                   = "896"
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([{
    name  = "api"
    image = "${var.api_image}:${var.api_image_tag}"



    portMappings = [{
      containerPort = 8080
      hostPort      = 0
      protocol      = "tcp"
    }]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = aws_cloudwatch_log_group.ecs.name
        "awslogs-region"        = var.aws_region
        "awslogs-stream-prefix" = "api"
      }
    }

    environment = [
      {
        name  = "ENV"
        value = var.environment
      },
      {
        name  = "PORT"
        value = "8080"
      },
      {
        name  = "MONGODB_NAME"
        value = var.mongodb_database_name
      },
      {
        name  = "CLIENT_ORIGIN"
        value = "https://${aws_lb.main.dns_name}"
      },
      {
        name  = "AWS_USER_POOL_ID"
        value = var.cognito_user_pool_id
      },
      {
        name  = "AWS_APP_CLIENT_ID"
        value = var.cognito_app_client_id
      },
      {
        name  = "AWS_S3_NAME"
        value = var.s3_name
      },
      {
        name  = "AWS_S3_URI"
        value = var.s3_uri
      },
      {
        name  = "AWS_REGION"
        value = var.aws_region
      },
      {
        name  = "AWS_SQS_QUEUE_NAME"
        value = var.sqs_queue_name
      },
      {
        name  = "AWS_SQS_QUEUE_URL"
        value = var.sqs_queue_url
      },
      {
        name = "FCM_PROJECT_ID"
        value = var.fcm_project_id
      }
    ]

    secrets = flatten([
        var.cognito_app_client_secret_arn != "" ? [{
        name      = "AWS_APP_CLIENT_SECRET"
        valueFrom = var.cognito_app_client_secret_arn
      }] : [],
        var.mongodb_uri_arn != "" ? [{
        name      = "MONGODB_URI"
        valueFrom = var.mongodb_uri_arn
      }] : [],
    ])

    memoryReservation = 768
  }])
}


resource "aws_ecs_task_definition" "client" {
  family                   = "${var.environment}-client"
  network_mode             = "bridge"
  requires_compatibilities = ["EC2"]
  cpu                      = "256"
  memory                   = "448"
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([{
    name  = "client"
    image = "${var.client_image}:${var.client_image_tag}"

    portMappings = [{
      containerPort = 80
      hostPort      = 0
      protocol      = "tcp"
    }]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = aws_cloudwatch_log_group.ecs.name
        "awslogs-region"        = var.aws_region
        "awslogs-stream-prefix" = "api"
      }
    }

    environment = [
      {
        name  = "__HALFTONE__API_URL"
        value = "https://${aws_lb.main.dns_name}"
      }
    ]

    memoryReservation = 384
  }])
}