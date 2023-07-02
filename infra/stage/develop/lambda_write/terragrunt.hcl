include "root" {
    path= find_in_parent_folders()
}

locals{
    env= read_terragrunt_config(find_in_parent_folders("env.hcl"))
}

terraform{
    source = "../../../modules/lambda-in-vpc"
}

dependency "ecr" {
    config_path = "../ecr"
    
    mock_outputs = {
        ecr_repository_url = "123456789012.dkr.ecr.ap-northeast-1.amazonaws.com"
    }
}
# VPCに依存
dependency "vpc" {
    config_path = "../vpc"

    mock_outputs = {
        vpc_id = "mock-vpc-id"
        private_subnet_ids = ["mock-private-subnet-id"]
        private_subnets_cidr_blocks = ["mock-private-subnet-cidr-block"]
        database_subnets_cidr_blocks = ["mock-database-subnet-cidr-block"]
    }
}
# RDSに依存
dependency "rds" {
    config_path = "../rds"

    mock_outputs = {
        rds_proxy_endpoint = "mock-rds-proxy-endpoint"
        rds_proxy_arn = "mock-rds-proxy-arn"
    }
}
inputs = {
    lambda_function_name = "feedays-cloud-write"
    lambda_function_description = "feedays-cloud-write-lambda-function"
    repo_url= dependency.ecr.outputs.ecr_repository_url
    image_tag= "write"
    memory_size = 128
    timeout = 30
    lambda_function_architecture = "arm64"

    # VPC設定
    subnet_ids = dependency.vpc.outputs.private_subnet_ids
    vpc_id = dependency.vpc.outputs.vpc_id

    # LamnbdaSGに必要
    vpc_private_subnets_cidr_blocks = dependency.vpc.outputs.private_subnets_cidr_blocks
    vpc_database_subnets_cidr_blocks = dependency.vpc.outputs.database_subnets_cidr_blocks

    # 環境変数はここで定義する
    variables = [
        rds_endpoint = dependency.rds.outputs.rds_proxy_endpoint,
        # RDSエンドポイント以外はEnv.hclから読み込むにした方が良い
        port = local.env.db_port,
        usename = local.env.db_username,
        db_name = local.env.db_name,
    ]

    managed_policy_arns = [
    # Lambda関数がCloudWatch Logsにログを書き込むための最低限の権限を提供します。
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
    # Lambda関数がVPC内のリソースにアクセスしながら実行するための最低限の権限（ネットワークインターフェースの作成、記述、削除、CloudWatch Logsへの書き込み権限）を提供します。
    "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole",
    # Amazon ECR に対する読み取り専用アクセスを付与
    "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
    # APIGW用
    "arn:aws:iam::aws:policy/AmazonAPIGatewayInvokeFullAccess",
    # RDS Proxy用
    "arn:aws:iam::aws:policy/AmazonRDSDataFullAccess",
    ]
}



