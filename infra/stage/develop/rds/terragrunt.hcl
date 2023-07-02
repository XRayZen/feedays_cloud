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
        vpc_cidr= "empty"
        vpc_private_subnets= ["subnet-123456", "subnet-123456"]
        vpc_private_subnets_cidr_blocks= ["empty", "empty"]
        vpc_database_subnet_group_name = "empty"
        vpc_database_subnets_cidr_blocks = ["empty", "empty"]
    }
}

inputs = {
    # ネットワーク設定
    vpc_id= dependency.vpc.outputs.vpc_id
    availability_zone = local.env.availability_zone
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
    # DBユーザー名
    db_username = local.env.db_username
    # DB名
    db_name = local.env.db_name
    tags = {
        Name = "${local.env.name}-db"
    }
}


