include "root" {
    path= find_in_parent_folders()
}

locals{
    env= read_terragrunt_config(find_in_parent_folders("env.hcl"))
}

terraform {
  source = "../../../modules/rds"
}

dependency "vpc" {
    config_path= "../vpc"

    mock_outputs= {
        vpc_id= "vpc-123456"
        vpc_cidr= "10.0.0.0/16"
        vpc_private_subnets= ["subnet-1234567", "subnet-12345678"]
        vpc_database_subnets= ["subnet-123456", "subnet-1234565"]
        vpc_database_subnet_group_name = "subnet-123456"
        vpc_private_subnets_cidr_blocks= ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
        vpc_database_subnets_cidr_blocks = ["10.0.110.0/24", "10.0.112.0/24", "10.0.113.0/24"]
    }
}

inputs = {
    # ネットワーク設定
    vpc_id= dependency.vpc.outputs.vpc_id
    availability_zone = local.env.locals.availability_zone
    vpc_cidr_block= dependency.vpc.outputs.vpc_cidr
    vpc_private_subnets = dependency.vpc.outputs.vpc_private_subnets
    vpc_database_subnets = dependency.vpc.outputs.vpc_database_subnets
    vpc_private_subnets_cidr_blocks = dependency.vpc.outputs.vpc_private_subnets_cidr_blocks
    vpc_database_subnet_group_name = dependency.vpc.outputs.vpc_database_subnet_group_name
    vpc_database_subnets_cidr_blocks = dependency.vpc.outputs.vpc_database_subnets_cidr_blocks
    # RDS設定
    db_engine = "mysql"
    db_engine_version = "8.0"
    db_parameter_group_family= "mysql8.0"
    db_instance_class = "db.t3.micro"
    db_storage_type = "gp2"
    db_allocated_storage = 20
    db_max_allocated_storage = 22
    # DB名
    db_name = local.env.locals.db_name
    # DBユーザー名
    db_username = local.env.locals.db_username
    # シークレット設定
    secret_version_stages = [local.env.locals.secret_stage]
    # Proxy設定
    rds_proxy_debug_logging = true
    rds_proxy_engine_family = "MYSQL"
    # RDS設定
    rds_create_cloudwatch_log_group = false
    rds_enabled_cloudwatch_logs_exports = ["general"]
    tags = {
        Name = "${local.env.locals.project_name}-db"
    }
}

