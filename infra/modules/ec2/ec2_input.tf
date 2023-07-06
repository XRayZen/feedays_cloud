# プロジェクト共通
variable "project_name" {
  type        = string
  description = "Project name"
}

variable "stage" {
  type        = string
  description = "Stage"
}

# EC2の設定
variable "key_name" {
  type = string
  description = "EC2インスタンスに関連付けられたキーペアの名前"
#   default = "${var.project_name}_${var.stage}_ec2_key"
}

variable "amazon_linux_ami_id" {
  type        = string
  description = "Amazon Linux AMI ID"
  # Amazon Linux 2023 AMI
  default = "ami-0cfc97bf81f2eadc4"
}

variable "ec2_instance_type" {
  type        = string
  description = "EC2 instance type"
  default     = "t2.micro"
}

variable "ec2_instance_name" {
  type        = string
  description = "EC2 instance name"
  default     = "ec2-instance"
}

# VPCの設定
variable "vpc_id" {
  type        = string
  description = "VPC ID"
}

variable "ec2_subnet_id" {
  type        = string
  description = "Subnet ID"
}

variable "vpc_public_subnets_cidr_blocks" {
  type        = list(string)
  description = "VPC public subnets CIDR blocks"
}

variable "availability_zone" {
  type        = string
  description = "Availability Zone"
}

# アクセス許可IPアドレス
variable "ec2_access_cidr_blocks" {
  type        = string
  description = "VPC access IP CIDR blocks(my public ip) OR (All Allow IP)"
  #   my public ipならば、以下のコマンドで取得できる
  #   curl -s ifconfig.me
  # Public IPアドレスを指定する場合は、以下のように指定する
  # "xx.xx.xx.xx/32"
  default = "0.0.0.0/0"
}

variable "rds_proxy_arn" {
  type        = string
  description = "RDS Proxy ARN"
}

# マネージドポリシーの設定
variable "managed_policy_arns" {
  type        = list(string)
  description = "Managed policy ARNs"
  default     = []
}

variable "tags" {
  type        = map(string)
  description = "Tags"
  default     = {}
}
