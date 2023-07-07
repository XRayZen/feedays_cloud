
module "ec2-instance" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "5.2.0"

  name = var.ec2_instance_name

  ami           = data.aws_ssm_parameter.amazonlinux_2023.value
  instance_type = var.ec2_instance_type
  key_name      = module.ec2_key_pair.key_pair_name
  monitoring    = true

  vpc_security_group_ids = [module.ec2_sg.security_group_id]
  availability_zone      = var.availability_zone
  subnet_id              = var.ec2_subnet_id

  # VPC内のインスタンスにパブリックIPアドレスを関連付けるかどうか
  associate_public_ip_address = true

  tags = var.tags
}
# Amazon Linux 3 の最新版AMIを取得
# Parameter Store のパブリックパラメーターを利用して AMI ID を取得
data "aws_ssm_parameter" "amazonlinux_2023" {
  name = "/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-6.1-x86_64" # x86_64
  # name = "/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-6.1-arm64" # ARM
  # name = "/aws/service/ami-amazon-linux-latest/al2023-ami-minimal-kernel-6.1-x86_64" # Minimal Image (x86_64)
  # name = "/aws/service/ami-amazon-linux-latest/al2023-ami-minimal-kernel-6.1-arm64" # Minimal Image (ARM)
}

# AMI 一覧から直接検索する方法
data "aws_ami" "amazonlinux_2023_x86" {
  most_recent = true
  owners = [ "amazon" ]
  filter {
    name = "name"

    values = [ "al2023-ami-*-kernel-6.1-x86_64" ] # x86_64
    # values = [ "al2023-ami-*-kernel-6.1-arm64" ] # ARM
    # values = [ "al2023-ami-minimal-*-kernel-6.1-x86_64" ] # Minimal Image (x86_64)
    # values = [ "al2023-ami-minimal-*-kernel-6.1-arm64" ] # Minimal Image (ARM)
  }
}

data "aws_ami" "amazonlinux_2023_arm" {
  most_recent = true
  owners = [ "amazon" ]
  filter {
    name = "name"

    # values = [ "al2023-ami-*-kernel-6.1-x86_64" ] # x86_64
    values = [ "al2023-ami-*-kernel-6.1-arm64" ] # ARM
    # values = [ "al2023-ami-minimal-*-kernel-6.1-x86_64" ] # Minimal Image (x86_64)
    # values = [ "al2023-ami-minimal-*-kernel-6.1-arm64" ] # Minimal Image (ARM)
  }
}



























