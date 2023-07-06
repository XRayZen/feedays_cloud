module "ec2_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.0"

  use_name_prefix = false
  name            = "${var.project_name}-${var.stage}-ec2-sg"
  description     = "Security group for ec2 instance"
  vpc_id          = var.vpc_id

  # インターネットからのアクセスを許可する
  # ルールは以下に記載されている
  # https://github.com/terraform-aws-modules/terraform-aws-security-group/blob/master/rules.tfb
  ingress_with_cidr_blocks = [
    # VPC内からのアクセスを許可する
    {
      description = "Allow HTTP from VPC"
      rule        = "http-80-tcp"
      cidr_blocks = join(",", var.vpc_public_subnets_cidr_blocks)
    },
    {
      description = "Allow SSH from My IP"
      rule        = "ssh-tcp"
      cidr_blocks = local.allowed_cidr
    },
  ]
  egress_with_cidr_blocks = [
    {
      # すべてのトラフィックを無差別に送信許可する
      description = "Allow all traffic"
      rule        = "all-all"
      cidr_blocks = "0.0.0.0/0"
    }
  ]

  tags = var.tags
}

# 自分のパブリックIP取得
data "http" "ifconfig" {
  url = "http://ipv4.icanhazip.com/"
}

variable "allowed_cidr" {
  default = null
}

locals {
  myip          = chomp(data.http.ifconfig.body)
  allowed_cidr  = (var.allowed_cidr == null) ? "${local.myip}/32" : var.allowed_cidr
}