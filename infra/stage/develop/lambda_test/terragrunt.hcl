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
        vpc_id = "vpc-1234567890abcdef0"
        vpc_private_subnets = ["subnet-1234567890abcdef0", "subnet-1234567890abcdef1", "subnet-1234567890abcdef2"]
        vpc_private_subnets_cidr_blocks = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
        vpc_database_subnets_cidr_blocks = ["10.0.110.0/24", "10.0.112.0/24", "10.0.113.0/24"]
    }
}
# RDSに依存
dependency "rds" {
    config_path = "../rds"

    mock_outputs = {
        rds_proxy_endpoint = "mock-rds-proxy-endpoint"
        rds_proxy_arn = "mock-rds-proxy-arn"
        db_password = "mock-db-password"
    }
}

inputs={
    lambda_function_name = "feedays-cloud-test"
    lambda_function_description = "feedays-cloud-test-lambda-function"
    repo_url= dependency.ecr.outputs.ecr_repository_url
    image_tag= "test"
    memory_size = 128
    timeout = 10
    lambda_function_architecture = "arm64"

    # VPC設定
    subnet_ids = dependency.vpc.outputs.vpc_private_subnets
    vpc_id = dependency.vpc.outputs.vpc_id

    # LamnbdaSGに必要
    vpc_private_subnets_cidr_blocks = dependency.vpc.outputs.vpc_private_subnets_cidr_blocks
    vpc_database_subnets_cidr_blocks = dependency.vpc.outputs.vpc_database_subnets_cidr_blocks
    rds_proxy_arn = dependency.rds.outputs.rds_proxy_arn
    # 環境変数はここで定義する
    variables = {
        region : local.env.locals.region,
        rds_endpoint: dependency.rds.outputs.rds_proxy_endpoint,
        # RDSエンドポイント以外はEnv.hclから読み込むにした方が良い
        db_port : local.env.locals.db_port,
        db_username : local.env.locals.db_username,
        db_name : local.env.locals.db_name,
        # パスワードはシークレットマネージャーから取得するので、使わない
        secret_stage : local.env.locals.secret_stage,
    }

    managed_policy_arns = [
    # Lambda関数がCloudWatch Logsにログを書き込むための最低限の権限を提供します。
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
    # Lambda関数がVPC内のリソースにアクセスしながら実行するための最低限の権限（ネットワークインターフェースの作成、記述、削除、CloudWatch Logsへの書き込み権限）を提供します。
    "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole",
    # Amazon ECR に対する読み取り専用アクセスを付与
    "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
    "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess",
    # APIGW用
    "arn:aws:iam::aws:policy/AmazonAPIGatewayInvokeFullAccess",
    # RDS Proxy用
    "arn:aws:iam::aws:policy/AmazonRDSDataFullAccess",
    # シークレットマネージャー用
    "arn:aws:iam::aws:policy/SecretsManagerReadWrite",
    ]
}


