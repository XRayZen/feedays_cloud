
# DBを作成
module "rds" {
  source  = "terraform-aws-modules/rds/aws"
  version = "6.0.0"
  # RDSインスタンスの名前
  identifier = "${var.db_name}-rds"

  db_name  = var.db_name
  username = var.db_username
  password = random_password.db-password.result
  port     = 3306

  iam_database_authentication_enabled = false

  engine               = var.db_engine
  engine_version       = var.db_engine_version
  family               = var.db_parameter_group_family
  major_engine_version = var.db_engine_version

  instance_class        = var.db_instance_class
  storage_type          = var.db_storage_type
  allocated_storage     = var.db_allocated_storage
  max_allocated_storage = var.db_max_allocated_storage
  # DBの変更を直ちに適用する
  apply_immediately = true

  multi_az               = false
  availability_zone      = var.availability_zone
  db_subnet_group_name   = var.vpc_database_subnet_group_name
  subnet_ids             = var.vpc_database_subnets
  vpc_security_group_ids = [module.rds_sg.security_group_id]

  maintenance_window         = "Mon:00:00-Mon:03:00"
  backup_window              = "03:00-06:00"
  backup_retention_period    = 0
  auto_minor_version_upgrade = false

  # 削除設定
  deletion_protection = false
  skip_final_snapshot = true

  tags = {
    Name = "${var.db_name}-rds"
  }
}
