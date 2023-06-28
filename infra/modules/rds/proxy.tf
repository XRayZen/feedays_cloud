
// RDS Proxy用のIAM
data "aws_iam_policy_document" "rds_proxy_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["rds.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "rds_proxy_role" {
  name               = "${var.db_name}-rds-proxy-role"
  assume_role_policy = data.aws_iam_policy_document.rds_proxy_assume_role.json
}

# RDS ProxyにSecrets Managerへのアクセス権限を付与
resource "aws_iam_role_policy" "rds_proxy_policy" {
  name   = "${var.db_name}-rds-proxy-policy"
  role   = aws_iam_role.rds_proxy_role.id
  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetResourcePolicy",
        "secretsmanager:GetSecretValue",
        "secretsmanager:DescribeSecret",
        "secretsmanager:ListSecretVersionIds"
      ],
      "Resource": "arn:aws:secretsmanager:*:*:*"
    }
  ]
}
POLICY
}
# RDS Proxy用のセキュリティグループ
resource "aws_security_group" "rds_proxy_security_group" {
  name        = "${var.db_name}-rds-proxy-security-group"
  description = "RDS Proxy security group"
  vpc_id      = var.vpc_id

  egress = {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# LambdaからRDS Proxyへのアクセス権限を付与
resource "aws_security_group_rule" "rds_proxy_security_group_rule" {
  security_group_id        = aws_security_group.rds_proxy_security_group.id
  type                     = "ingress"
  from_port                = 3306
  to_port                  = 3306
  protocol                 = "tcp"
}

# RDS Proxy
resource "aws_db_proxy" "rds_proxy" {
  name                   = "${var.db_name}-rds-proxy"
  engine_family          = "MYSQL"
  role_arn               = aws_iam_role.rds_proxy_role.arn
  vpc_security_group_ids = [aws_security_group.rds_proxy_security_group.id]
  vpc_subnet_ids         = var.subnet_ids

  auth {
    secret_arn = aws_secretsmanager_secret.db-password.arn
    username   = var.db_username
  }
}

# ターゲットグループは、プロキシが接続できるデータベースのコレクションです。
# 現在、各ターゲットグループを 1 つの RDS DB インスタンスまたは Aurora DB クラスターに関連付けることができます。
resource "aws_db_proxy_default_target_group" "rds_proxy_default_target_group" {
  db_proxy_name = aws_db_proxy.rds_proxy.name
}

resource "aws_db_proxy_target" "example" {
  db_instance_identifier = aws_db_instance.mysql_db_instance.id
  db_proxy_name          = aws_db_proxy.rds_proxy.name
  target_group_name      = "default"
}





