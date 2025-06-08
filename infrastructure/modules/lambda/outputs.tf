output "lambda_function_arn" {
  description = "The ARN of the Lambda function"
  value       = aws_lambda_function.photo_processor.arn
}

output "lambda_function_role" {
  description = "The IAM role assumed by the Lambda function"
  value       = aws_lambda_function.photo_processor.role
}

