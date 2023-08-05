# RDS Proxy
module "rds_proxy" {
  source  = "terraform-aws-modules/rds-proxy/aws"
  version = "3.0.0"

  name                   = "${var.db_name}-rds-proxy"
  iam_role_name          = "${var.db_name}-rds-proxy-role"
  vpc_subnet_ids         = var.vpc_private_subnets
  vpc_security_group_ids = [module.rds_proxy_sg.security_group_id]

  engine_family = var.rds_proxy_engine_family
  debug_logging = var.rds_proxy_debug_logging
  # TLSを有効にすると、GORMでRDS Proxyのエンドポイントに接続できないので無効化しておく
  require_tls = false

  target_db_instance     = true
  db_instance_identifier = module.rds.db_instance_identifier

  endpoints = {
    read_write = {
      name                   = "rw-endpoint"
      vpc_subnet_ids         = var.vpc_private_subnets
      vpc_security_group_ids = [module.rds_proxy_sg.security_group_id]
      tags = {
        Name = "mysql-rw-endpoint"
      }
    },
    read_only = {
      name                   = "ro-endpoint"
      vpc_subnet_ids         = var.vpc_private_subnets
      vpc_security_group_ids = [module.rds_proxy_sg.security_group_id]
      tagete_role            = "READ_ONLY"
      tags = {
        Name = "mysql-ro-endpoint"
      }
    }
  }

  auth = {
    (var.db_username) = {
      secret_arn  = aws_secretsmanager_secret.superuser.arn
      description = aws_secretsmanager_secret.superuser.description
    }
  }

  # セッション固定フィルタを設定して接続を再利用できる様にする
  session_pinning_filters = [
    "EXCLUDE_VARIABLE_SETS",
    # "INCLUDE_ALL"
  ]

  tags = {
    Name = "rds-proxy"
  }
}
