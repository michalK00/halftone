# environments /dev/main

terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1"
    }
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Environment = var.environment
      Project     = var.project_name
      ManagedBy   = "Terraform"
    }
  }
}

data "aws_caller_identity" "current" {}


module "networking" {
  source = "../../modules/networking"

  environment        = var.environment
  vpc_cidr          = var.vpc_cidr
  availability_zones = var.availability_zones
}

module "storage" {
  source = "../../modules/storage"

  project_name    = var.project_name
  environment     = var.environment
  allowed_origins = ["*"] # Allow all origins since there is no domain
}

# Authentication
module "auth" {
  source = "../../modules/auth"

  project_name = var.project_name
  environment  = var.environment
}

module "ecr" {
  source = "../../modules/ecr"

  project_name = var.project_name
  environment  = var.environment
}


# Security
module "security" {
  source = "../../modules/security"

  project_name      = var.project_name
  environment       = var.environment
  photos_bucket_arn = module.storage.photos_bucket_arn
  logs_bucket_arn   = module.storage.logs_bucket_arn
  user_pool_arn     = module.auth.user_pool_arn
}

module "database" {
  source = "../../modules/database"

  environment            = var.environment
  project_name           = var.project_name
  vpc_id                 = module.networking.vpc_id
  subnet_ids             = module.networking.private_subnet_ids
  app_security_group_ids = [module.networking.ecs_tasks_security_group_id, module.networking.database_security_group_id]

  master_username = var.docdb_master_username
  master_password = var.docdb_master_password

  instance_count = 1
  instance_type = "db.t3.medium"
}

# MongoDB URI Secret
resource "aws_secretsmanager_secret" "mongodb_uri" {
  name = "${var.environment}-db-uri"

  recovery_window_in_days = 0
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_secretsmanager_secret_version" "mongodb_uri" {
  secret_id = aws_secretsmanager_secret.mongodb_uri.id
  secret_string = "mongodb://${var.docdb_master_username}:${var.docdb_master_password}@${module.database.cluster_endpoint}:27017/${var.mongodb_database_name}?authSource=admin&tls=true&tlsCAFile=/opt/global-bundle.pem&replicaSet=rs0&readPreference=secondaryPreferred&retryWrites=false"
}

module "notifications" {
  source = "../../modules/notifications"

  project_name      = var.project_name
  environment       = var.environment

  sns_topic_name = "alert-notifications"

  email_subscriber = var.notification_email
  sns_display_name = "alert-notifications-${var.environment}"
}

module "lambda" {
  source = "../../modules/lambda"

  project_name      = var.project_name
  environment       = var.environment
  function_name     = "photo-uploads-handler"

  tags = {
    Environment = var.environment
    Project     = var.project_name
    Module      = "lambda"
  }
  sns_topic_arn      = module.notifications.sns_topic_arn
  sqs_queue_arn      = module.messaging.sqs_queue_arn
  photos_bucket_arn  = module.storage.photos_bucket_arn
  source_dir         = var.lambda_source_dir
}

module "messaging" {
  source = "../../modules/messaging"

  aws_sns_topic_arn   = module.notifications.sns_topic_arn
  environment         = var.environment
  project_name        = var.project_name
  sqs_queue_name      = "photo-uploads-queue"
}

# ECS Module
module "ecs" {
  source = "../../modules/ecs"

  environment                 = var.environment
  aws_region                  = var.aws_region
  vpc_id                      = module.networking.vpc_id
  public_subnet_ids           = module.networking.public_subnet_ids
  private_subnet_ids          = module.networking.private_subnet_ids
  alb_security_group_id       = module.networking.alb_security_group_id
  ecs_tasks_security_group_id = module.networking.ecs_tasks_security_group_id

  # Container configuration
  api_image        = module.ecr.api_repository_url
  api_image_tag    = var.api_image_tag
  client_image     = module.ecr.client_repository_url
  client_image_tag = var.client_image_tag

  # Cognito
  cognito_user_pool_id          = module.auth.user_pool_id
  cognito_app_client_id         = module.auth.user_pool_client_id
  cognito_app_client_secret_arn = module.auth.client_secret_arn

  # S3
  s3_name = module.storage.photos_bucket_name
  s3_uri  = "s3://${module.storage.photos_bucket_name}"

  # MongoDB
  mongodb_uri_arn = aws_secretsmanager_secret.mongodb_uri.arn
  mongodb_database_name = var.mongodb_database_name

  # Service configuration
  api_desired_count      = var.api_desired_count
  frontend_desired_count = var.frontend_desired_count

  # SQS
  sqs_queue_name = module.messaging.sqs_queue_name
  sqs_queue_url  = module.messaging.sqs_queue_url
  sqs_queue_arn  = module.messaging.sqs_queue_arn
}

# Outputs
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

output "ecr_repositories" {
  description = "ECR repository URLs"
  value = {
    api    = module.ecr.api_repository_url
    client = module.ecr.client_repository_url
  }
}

output "ecr_registry" {
  description = "ECR registry URL"
  value       = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.aws_region}.amazonaws.com"
}