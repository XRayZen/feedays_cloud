
module "rds_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.0"

  name        = "${var.db_name}-rds-sg"
  description = "RDS security group"
  vpc_id      = var.vpc_id

  revoke_rules_on_delete = true

  # プライベートサブネットからのアクセスを許可する
  ingress_with_cidr_blocks = [
    {
      description = "private subnet access"
      rule        = "mysql-tcp"
      cidr_blocks = join(",", var.vpc_private_subnets_cidr_blocks)
    },
  ]

  tags = {
    Name = "${var.db_name}-rds-sg"
  }
}

module "rds_proxy_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.0"

  name        = "${var.db_name}-rds-proxy-sg"
  description = "RDS Proxy security group"
  vpc_id      = var.vpc_id
  # ルール自体を削除する前に、イングレスとイグレスのルールにアタッチされたすべてのSecurity Groupを取り消すようTerraformに指示する。EMRを有効にする。
  revoke_rules_on_delete = true

  # イングレスルール
  ingress_with_cidr_blocks = [
    #  Private Subnetからのアクセスを許可する
    {
      description = "private subnet access"
      rule        = "mysql-tcp"
      cidr_blocks = join(",", var.vpc_private_subnets_cidr_blocks)
    },
  ]

  egress_with_cidr_blocks = [
    {
      # データベースサブネットへのアウトバンドを許可する
      description= "Database subnet access"
      rule        = "mysql-tcp"
      cidr_blocks = join(",", var.vpc_database_subnets_cidr_blocks)
    }
  ]

  tags = {
    Name = "${var.db_name}-rds-proxy-sg"
  }
}
