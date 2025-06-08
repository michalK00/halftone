resource "aws_cognito_user_pool" "main" {
  name = "${var.project_name}-${var.environment}-users"

  username_attributes      = ["email"]
  auto_verified_attributes = ["email"]

  password_policy {
    minimum_length    = 8
    require_lowercase = true
    require_numbers   = true
    require_symbols   = true
    require_uppercase = true
  }

  email_configuration {
    email_sending_account = "COGNITO_DEFAULT"
  }

  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }

  mfa_configuration = "OFF"

  tags = {
    Name        = "${var.project_name}-${var.environment}-user-pool"
    Environment = var.environment
  }
}

resource "aws_cognito_user_pool_client" "main" {
  name         = "${var.project_name}-${var.environment}-client"
  user_pool_id = aws_cognito_user_pool.main.id

  explicit_auth_flows = [
    "ALLOW_USER_PASSWORD_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH",
    "ALLOW_USER_SRP_AUTH"
  ]

  access_token_validity  = 1
  id_token_validity      = 1
  refresh_token_validity = 5

  token_validity_units {
    access_token  = "hours"
    id_token      = "hours"
    refresh_token = "days"
  }

  prevent_user_existence_errors = "ENABLED"
  enable_token_revocation      = true

  read_attributes = [
    "email",
    "email_verified",
    "sub"
  ]

  generate_secret = true
}

resource "aws_secretsmanager_secret" "client_secret" {
  name = "${var.project_name}-${var.environment}-client-secret"

  recovery_window_in_days = 0
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_secretsmanager_secret_version" "client_secret" {
  secret_id     = aws_secretsmanager_secret.client_secret.id
  secret_string = aws_cognito_user_pool_client.main.client_secret
}