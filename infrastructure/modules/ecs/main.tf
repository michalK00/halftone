# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "${var.environment}-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Name        = "${var.environment}-cluster"
    Environment = var.environment
  }
}

resource "aws_cloudwatch_log_group" "ecs" {
  name              = "/ecs/${var.environment}"
  retention_in_days = 1

  tags = {
    Name        = "/ecs/${var.environment}"
    Environment = var.environment
  }
}

data "aws_ami" "ecs_optimized" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-ecs-hvm-*-x86_64-ebs"]
  }
}

resource "aws_iam_role" "ecs_instance" {
  name = "${var.environment}-ecs-instance-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_instance" {
  role       = aws_iam_role.ecs_instance.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"
}

resource "aws_iam_instance_profile" "ecs_instance" {
  name = "${var.environment}-ecs-instance-profile"
  role = aws_iam_role.ecs_instance.name
}

# Task execution role
resource "aws_iam_role" "ecs_task_execution" {
  name = "${var.environment}-ecs-task-execution"

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
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# Policy for secrets
resource "aws_iam_policy" "ecs_secrets" {
  name = "${var.environment}-ecs-secrets"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "secretsmanager:GetSecretValue"
      ]
      Resource = compact([
        var.cognito_app_client_secret_arn,
        var.mongodb_uri_arn
      ])
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_secrets" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = aws_iam_policy.ecs_secrets.arn
}

# Task role
resource "aws_iam_role" "ecs_task" {
  name = "${var.environment}-ecs-task"

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
}

# S3 access for tasks
resource "aws_iam_policy" "ecs_s3" {
  name = "${var.environment}-ecs-s3"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject"
      ]
      Resource = [
        "arn:aws:s3:::${var.s3_name}/*"
      ]
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_s3" {
  role       = aws_iam_role.ecs_task.name
  policy_arn = aws_iam_policy.ecs_s3.arn
}

# ECS Instances
resource "aws_instance" "ecs" {
  count                  = var.instance_count
  ami                    = data.aws_ami.ecs_optimized.id
  instance_type          = var.instance_type
  subnet_id              = var.private_subnet_ids[count.index % length(var.private_subnet_ids)]
  vpc_security_group_ids = [var.ecs_tasks_security_group_id]
  iam_instance_profile   = aws_iam_instance_profile.ecs_instance.name

  associate_public_ip_address = true

  user_data = <<-EOF
            #!/bin/bash
            echo "ECS_CLUSTER=${aws_ecs_cluster.main.name}" >> /etc/ecs/ecs.config
            EOF

  tags = {
    Name        = "${var.environment}-ecs-${count.index + 1}"
    Environment = var.environment
  }
}

resource "tls_private_key" "key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "cert" {
  private_key_pem = tls_private_key.key.private_key_pem

  subject {
    common_name  = aws_lb.main.dns_name
    organization = "Example Organization"
  }

  validity_period_hours = 8760 # Valid for 1 year

  dns_names = [
    aws_lb.main.dns_name
  ]

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]

  depends_on = [aws_lb.main]
}

resource "aws_acm_certificate" "self_signed" {
  private_key      = tls_private_key.key.private_key_pem
  certificate_body = tls_self_signed_cert.cert.cert_pem

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name = "self-signed-certificate"
  }
}

# Application Load Balancer
resource "aws_lb" "main" {
  name               = "${var.environment}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [var.alb_security_group_id]
  subnets            = var.public_subnet_ids

  enable_deletion_protection = false

  tags = {
    Name        = "${var.environment}-alb"
    Environment = var.environment
  }
}

resource "aws_lb_target_group" "api" {
  name        = "${var.environment}-api-tg"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "instance"

  health_check {
    enabled             = true
    healthy_threshold   = 2
    interval            = 30
    matcher             = "200"
    path                = "/health"
    port                = "traffic-port"
    protocol            = "HTTP"
    timeout             = 5
    unhealthy_threshold = 2
  }
}

resource "aws_lb_target_group" "client" {
  name        = "${var.environment}-client-tg"
  port        = 3000
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "instance"

  health_check {
    enabled             = true
    healthy_threshold   = 2
    interval            = 30
    matcher             = "200"
    path                = "/"
    port                = "traffic-port"
    protocol            = "HTTP"
    timeout             = 5
    unhealthy_threshold = 2
  }
}

# ALB Listener
resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.main.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS-1-2-2017-01"
  certificate_arn   = aws_acm_certificate.self_signed.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.client.arn
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.main.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}
# Listener Rules
resource "aws_lb_listener_rule" "api_https" {
  listener_arn = aws_lb_listener.https.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api.arn
  }

  condition {
    path_pattern {
      values = ["/api/*"]
    }
  }
}

resource "aws_lb_listener_rule" "api" {
  listener_arn = aws_lb_listener.http.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api.arn
  }

  condition {
    path_pattern {
      values = ["/api/*"]
    }
  }
}

# Task Definitions
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
      }] : []
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
        value = "https://${aws_lb.main.dns_name}/api"
      }
    ]

    memoryReservation = 384
  }])
}

# ECS Services
resource "aws_ecs_service" "api" {
  name            = "${var.environment}-api"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.api.arn
  desired_count   = var.api_desired_count
  launch_type     = "EC2"

  load_balancer {
    target_group_arn = aws_lb_target_group.api.arn
    container_name   = "api"
    container_port   = 8080
  }

  depends_on = [aws_lb_listener.http]
}

resource "aws_ecs_service" "client" {
  name            = "${var.environment}-client"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.client.arn
  desired_count   = var.frontend_desired_count
  launch_type     = "EC2"

  load_balancer {
    target_group_arn = aws_lb_target_group.client.arn
    container_name   = "client"
    container_port   = 80
  }

  depends_on = [aws_lb_listener.http]
}