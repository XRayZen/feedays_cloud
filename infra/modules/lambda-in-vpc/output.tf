output "lambda_function_name" {
  value = aws_lambda_function.lambda_func_.function_name
}

output "lambda_function_invoke_arn" {
  value = aws_lambda_function.lambda_func_.invoke_arn
}

output "lambda_security_group_id" {
  value = module.lambda_sg.security_group_id
}