
module "ec2-instance" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "5.2.0"

  name = var.ec2_instance_name

  ami           = var.amazon_linux_ami_id
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
data "aws_ssm_parameter" "amzn2_latest_ami" {
    name = "/aws/service/ami-amazon-linux-latest/amzn2-ami-hvm-x86_64-gp2"
}
