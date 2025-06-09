resource "aws_codedeploy_app" "api" {
  compute_platform = "ECS"
  name             = "${var.environment}-api-codedeploy"

  tags = {
    Name        = "${var.environment}-api-codedeploy"
    Environment = var.environment
  }
}

resource "aws_codedeploy_deployment_group" "api_with_alarms" {
  app_name               = aws_codedeploy_app.api.name
  deployment_group_name  = "${var.environment}-api-deployment-group-with-alarms"
  service_role_arn      = aws_iam_role.codedeploy.arn
  deployment_config_name = "CodeDeployDefault.ECSAllAtOnceBlueGreen"

  auto_rollback_configuration {
    enabled = true
    events  = ["DEPLOYMENT_FAILURE", "DEPLOYMENT_STOP_ON_ALARM", "DEPLOYMENT_STOP_ON_INSTANCE_FAILURE"]
  }

  alarm_configuration {
    alarms = [
      aws_cloudwatch_metric_alarm.api_high_error_rate.alarm_name,
      aws_cloudwatch_metric_alarm.api_low_healthy_hosts.alarm_name
    ]
    enabled = true
  }

  blue_green_deployment_config {
    terminate_blue_instances_on_deployment_success {
      action                         = "TERMINATE"
      termination_wait_time_in_minutes = 5
    }

    deployment_ready_option {
      action_on_timeout = "CONTINUE_DEPLOYMENT"
    }

    green_fleet_provisioning_option {
      action = "COPY_AUTO_SCALING_GROUP"
    }
  }

  ecs_service {
    cluster_name = aws_ecs_cluster.main.name
    service_name = aws_ecs_service.api.name
  }

  load_balancer_info {
    target_group_info {
      name = aws_lb_target_group.api.name
    }
  }

  tags = {
    Name        = "${var.environment}-api-deployment-group-with-alarms"
    Environment = var.environment
  }
}