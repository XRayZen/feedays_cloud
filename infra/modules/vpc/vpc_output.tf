output "vpc_id" {
  value       = module.vpc.vpc_id
  description = "VPC ID"
}

output "vpc_cidr" {
  value       = module.vpc.vpc_cidr_block
  description = "VPC CIDR"
}

output "public_subnet_id_1" {
  value       = data.aws_subnet.public[0].id
  description = "Public Subnet ID 1"
  type        = string
}

output "public_subnet_id_2" {
  value       = data.aws_subnet.public[1].id
  description = "Public Subnet ID 2"
  type        = string
}

output "private_subnet_id_1" {
  value       = data.aws_subnet.private[0].id
  description = "Private Subnet ID 1"
  type        = string
}

output "private_subnet_id_2" {
  value       = data.aws_subnet.private[1].id
  description = "Private Subnet ID 2"
  type        = string
}

