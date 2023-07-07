
module "lambda_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.0"

  name        = "lambda_sg"
  description = "Security group for lambda function"
  vpc_id      = var.vpc_id

  #   インターネットへのアクセスを許可する
  ingress_with_cidr_blocks = [
    {
      description = "Allow HTTP from VPC"
      # ルールは以下に記載されている
      # https://github.com/terraform-aws-modules/terraform-aws-security-group/blob/master/rules.tf
      rule        = "http-80-tcp"
      cidr_blocks = join(",", var.vpc_private_subnets_cidr_blocks)
    },
  ]

  # egress_with_cidr_blocks = [
  #   {
  #     # すべてのトラフィックをデータベースサブネットに送信許可する
  #     description = "Allow all traffic"
  #     rule        = "all-all"
  #     cidr_blocks = join(",", var.vpc_database_subnets_cidr_blocks)
  #   },
  #   {
  #     # すべてのトラフィックを無差別に送信許可する
  #     description = "Allow all traffic"
  #     rule        = "all-all"
  #     cidr_blocks = "0.0.0.0/0"
  #   }
  # ]
  tags = {
    Name        = "lambda_sg"
    description = "Security group for lambda function"
  }
}

