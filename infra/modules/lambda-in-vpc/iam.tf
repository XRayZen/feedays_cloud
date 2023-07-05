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
# RDS Proxyへのアクセス権限を付与
data "aws_iam_policy_document" "lambda_rds_proxy_policy_document" {
  statement {
    sid       = ""
    effect    = "Allow"
    actions   = ["rds-db:connect"]
    resources = [var.rds_proxy_arn]
  }
}

resource "aws_iam_policy" "lambda_rds_proxy_policy" {
  name   = "${var.lambda_function_name}-rds-proxy-policy"
  policy = data.aws_iam_policy_document.lambda_rds_proxy_policy_document.json
}

resource "aws_iam_role_policy_attachment" "lambda_rds_proxy_policy_attachment" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.lambda_rds_proxy_policy.arn
}

# シークレットマネージャーへのアクセス権限を付与
data "aws_iam_policy_document" "lambda_secrets_manager_policy_document" {
  statement {
    sid    = ""
    effect = "Allow"
    actions = [
      "secretsmanager:*",
      "kms:DescribeKey",
      "kms:ListAliases",
      "kms:ListKeys",
      "kms:Decrypt",
      "lambda:ListFunctions",
      "rds:DescribeDBClusters",
      "rds:DescribeDBInstances",
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "lambda_secrets_manager_policy" {
  name   = "${var.lambda_function_name}-secrets-manager-policy"
  policy = data.aws_iam_policy_document.lambda_secrets_manager_policy_document.json
}

resource "aws_iam_role_policy_attachment" "lambda_secrets_manager_policy_attachment" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.lambda_secrets_manager_policy.arn
}

