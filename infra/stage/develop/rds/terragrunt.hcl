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
    }
}

dependency "lambda_read"{
    config_path= "../lambda_read"

    mock_outputs= {
        lambda_arn= "arn:aws:lambda:us-east-1:123456789012:function:lambda-read"

    }
}

inputs = {
    # ネットワーク設定
    vpc_id= dependency.vpc.outputs.vpc_id
    availability_zone = "us-east-1a"
    vpc_cidr_block= dependency.vpc.outputs.vpc_cidr
    subnet_ids= [dependency.vpc.outputs.private_subnet_id_1, dependency.vpc.outputs.private_subnet_id_2]
    # RDSには複数のLambdaからアクセスするため、セキュリティグループを分離して定義する必要がある

    # RDS設定
}


