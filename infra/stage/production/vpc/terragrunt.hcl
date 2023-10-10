include "root" {
    path= find_in_parent_folders()
}

locals{
    env= read_terragrunt_config(find_in_parent_folders("env.hcl"))
}

terraform{
    source="../../../modules/vpc"
}

inputs ={
    project_name= local.env.locals.project_name

    cidr = "10.0.0.0/16"
    availability_zones = ["us-east-1a", "us-east-1b", "us-east-1c"]
    private_subnets =["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
    public_subnets =["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
    database_subnets =["10.0.110.0/24", "10.0.112.0/24", "10.0.113.0/24"]

    # LambdaがVPC内のリソースとにアクセスするためには、NAT Gatewayが必要
    enable_nat_gateway = true
    # single_nat_gateway = trueの場合、すべてのプライベート・サブネットはこの単一のNATゲートウェイを経由してインターネット・トラフィックをルーティングします。
    # NATゲートウェイはpublic_subnetsブロックの最初のパブリック・サブネットに配置されます。
    single_nat_gateway = true
    one_nat_gateway_per_az = false
}

