

resource "aws_lambda_function" "lambda_func_" {
  function_name = var.lambda_function_name
  description   = var.lambda_function_description
  memory_size   = var.memory_size
  timeout       = var.timeout
  role          = aws_iam_role.lambda_role.arn
  architectures = [var.lambda_function_architecture]

  package_type = "Image"
  image_uri    = "${var.repo_url}:${var.image_tag}"

  vpc_config {
    subnet_ids         = var.subnet_ids
    security_group_ids = var.security_group_ids
  }
  # environment {
  #   variables = var.environment_variables
  # }
}

resource "aws_iam_role" "lambda_role" {
  name                = "${var.lambda_function_name}-role"
  assume_role_policy  = data.aws_iam_policy_document.lambda_role_policy.json
  managed_policy_arns = var.managed_policy_arns
}

data "aws_iam_policy_document" "lambda_role_policy" {
  statement {
    sid     = ""
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}
