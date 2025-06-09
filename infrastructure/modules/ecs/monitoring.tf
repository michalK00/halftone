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

resource "aws_cloudwatch_metric_alarm" "api_high_error_rate" {
  alarm_name          = "${var.environment}-api-high-error-rate"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "HTTPCode_Target_5XX_Count"
  namespace           = "AWS/ApplicationELB"
  period              = "60"
  statistic           = "Sum"
  threshold           = "5"
  alarm_description   = "High error rate detected"
  treat_missing_data  = "notBreaching"

  dimensions = {
    TargetGroup = aws_lb_target_group.api.arn_suffix
  }

  tags = {
    Name        = "${var.environment}-api-high-error-rate"
    Environment = var.environment
  }
}

resource "aws_cloudwatch_metric_alarm" "api_low_healthy_hosts" {
  alarm_name          = "${var.environment}-api-low-healthy-hosts"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "HealthyHostCount"
  namespace           = "AWS/ApplicationELB"
  period              = "60"
  statistic           = "Average"
  threshold           = "1"
  alarm_description   = "Low healthy host count"
  treat_missing_data  = "breaching"

  dimensions = {
    TargetGroup = aws_lb_target_group.api.arn_suffix
  }

  tags = {
    Name        = "${var.environment}-api-low-healthy-hosts"
    Environment = var.environment
  }
}