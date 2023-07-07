include "root" {
    path= find_in_parent_folders()
}

locals{
    env= read_terragrunt_config(find_in_parent_folders("env.hcl"))
}

terraform {
  source = "../../../modules/ec2"
}

dependency "vpc" {
    config_path= "../vpc"

    mock_outputs= {
        vpc_id= "vpc-123456"
        vpc_public_subnet_1_id = "subnet-123456"
        vpc_public_subnets_cidr_blocks = ["10.0.110.0/24", "10.0.112.0/24", "10.0.113.0/24"]
    }
}

dependency "rds"{
    config_path= "../rds"
    mock_outputs ={
        rds_proxy_arn = "arn:aws:rds:ap-northeast-1:123456789012:db:my-postgres-db"
    }
}

inputs ={
    project_name= local.env.locals.project_name
    stage = local.env.locals.stage
    # EC2の設定
    key_name= "${local.env.locals.project_name}_${local.env.locals.stage}_ec2_key"
    amazon_linux_ami_id = "ami-0cfc97bf81f2eadc4"
    ec2_instance_type= "t2.micro"
    ec2_instance_name= "ec2-instance"
    # VPCの設定
    vpc_id= dependency.vpc.outputs.vpc_id
    ec2_subnet_id= dependency.vpc.outputs.vpc_public_subnet_1_id
    vpc_public_subnets_cidr_blocks = dependency.vpc.outputs.vpc_public_subnets_cidr_blocks
    availability_zone= local.env.locals.availability_zone
    # 今のところはEC2にアクセスするIPを指定しない
    rds_proxy_arn = dependency.rds.outputs.rds_proxy_arn
    # 付与するマネージドポリシーのARN
    managed_policy_arns =[
        "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore",
        "arn:aws:iam::aws:policy/AmazonSSMFullAccess",
        "arn:aws:iam::aws:policy/AmazonRDSDataFullAccess",
        # VPCへの権限を付与
        "arn:aws:iam::aws:policy/AmazonEC2FullAccess",
    ]
}























