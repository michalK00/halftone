resource "aws_codedeploy_app" "api" {
  compute_platform = "ECS"
  name             = "${var.environment}-api-codedeploy"

  tags = {
    Name        = "${var.environment}-api-codedeploy"
    Environment = var.environment
  }
}

resource "aws_codedeploy_deployment_group" "api" {
  app_name               = aws_codedeploy_app.api.name
  deployment_config_name = "CodeDeployDefault.ECSAllAtOnce"
  deployment_group_name  = "${var.environment}-api-deployment-group"
  service_role_arn      = aws_iam_role.codedeploy.arn

  deployment_style {
    deployment_type = "BLUE_GREEN"
    deployment_option = "WITH_TRAFFIC_CONTROL"
  }

  auto_rollback_configuration {
    enabled = true
    events  = ["DEPLOYMENT_FAILURE"]
  }

  blue_green_deployment_config {
    terminate_blue_instances_on_deployment_success {
      action                         = "TERMINATE"
      termination_wait_time_in_minutes = 5
    }

    deployment_ready_option {
      action_on_timeout = "CONTINUE_DEPLOYMENT"
    }
  }

  ecs_service {
    cluster_name = aws_ecs_cluster.main.name
    service_name = aws_ecs_service.api.name
  }

  load_balancer_info {
    target_group_pair_info {
      prod_traffic_route {
        listener_arns = [aws_lb_listener.https.arn]
      }
      target_group {
        name = aws_lb_target_group.api.name
      }
      target_group {
        name = aws_lb_target_group.api_green.name
      }
    }
  }


  tags = {
    Name        = "${var.environment}-api-deployment-group-with-alarms"
    Environment = var.environment
  }
}