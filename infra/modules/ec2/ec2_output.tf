output "key_name" {
  description = "EC2インスタンスに関連付けられたキーペアの名前"
  value       = module.ec2_key_pair.key_pair_name
}

output "ec2_instance_ip" {
  description = "インスタンスに割り当てられたパブリックIPアドレス（該当する場合）。注: インスタンスで aws_eip を使用している場合は、EIP のアドレスを直接参照し、public_ip を使用しないでください。"
  value       = module.ec2-instance.public_ip
}
