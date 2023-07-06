# モジュールIAMは複雑なので使わない

resource "aws_iam_role" "ec2_role" {
  name                = "${var.project_name}_${var.stage}_ec2_role"
  assume_role_policy  = data.aws_iam_policy_document.ec2_role_policy.json
  description         = "EC2 Role"
  managed_policy_arns = var.managed_policy_arns
}

data "aws_iam_policy_document" "ec2_role_policy" {
  statement {
    sid     = ""
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

# RDS Proxyへのアクセス権限を付与
data "aws_iam_policy_document" "ec2_rds_proxy_policy_document" {
  statement {
    sid       = ""
    effect    = "Allow"
    actions   = ["rds-db:connect"]
    resources = [var.rds_proxy_arn]
  }
}

resource "aws_iam_policy" "ec2_rds_proxy_policy" {
  name   = "${var.project_name}_${var.stage}_ec2_rds_proxy_policy"
  policy = data.aws_iam_policy_document.ec2_rds_proxy_policy_document.json
}

# EC2のIAMロールにRDS Proxyへのアクセス権限を付与
resource "aws_iam_role_policy_attachment" "ec2_rds_proxy_policy_attachment" {
  role       = aws_iam_role.ec2_role.name
  policy_arn = aws_iam_policy.ec2_rds_proxy_policy.arn
}



