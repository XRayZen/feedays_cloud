# RDS Proxy
module "rds_proxy" {
  source  = "terraform-aws-modules/rds-proxy/aws"
  version = "3.0.0"

  name                   = "${var.db_name}-rds-proxy"
  iam_role_name          = "${var.db_name}-rds-proxy-role"
  vpc_subnet_ids         = var.vpc_private_subnets
  vpc_security_group_ids = [module.rds_proxy_sg.security_group_id]

  engine_family = "MYSQL"
  debug_logging = true

  # Target RDS instance
  db_instance_identifier = module.rds.db_instance_identifier

  endpoints = {
    read_write = {
      name                   = "mysql-rw-endpoint"
      vpc_subnet_ids         = var.vpc_private_subnets
      vpc_security_group_ids = [module.rds_proxy_sg.security_group_id]
      tags = {
        Name = "mysql-rw-endpoint"
      }
    },
    read_only = {
      name                   = "mysql-ro-endpoint"
      vpc_subnet_ids         = var.vpc_private_subnets
      vpc_security_group_ids = [module.rds_proxy_sg.security_group_id]
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

  tags = {
    Name = "rds-proxy"
  }
}
