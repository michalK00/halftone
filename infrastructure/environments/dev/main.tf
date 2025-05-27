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

# Networking
module "networking" {
  source = "../../modules/networking"

  environment        = var.environment
  vpc_cidr          = var.vpc_cidr
  availability_zones = var.availability_zones
}

# Storage
module "storage" {
  source = "../../modules/storage"

  project_name    = var.project_name
  environment     = var.environment
  allowed_origins = var.allowed_origins
}

# Authentication
module "auth" {
  source = "../../modules/auth"

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

# DocumentDB Preparation
# module "documentdb" {
#   source = "../../modules/database"
#
#   environment           = var.environment
#   vpc_id               = module.networking.vpc_id
#   subnet_ids           = module.networking.private_subnet_ids
#   security_group_id    = module.networking.database_security_group_id
#   master_username      = var.docdb_master_username
#   master_password      = var.docdb_master_password
# }
