output "vpc_id" {
  value       = module.vpc.vpc_id
  description = "VPC ID"
}

output "vpc_cidr" {
  value       = module.vpc.vpc_cidr_block
  description = "VPC CIDR"
}

output "vpc_private_subnets" {
  value       = module.vpc.private_subnets
  description = "VPC Private Subnets IDs"
}

output "vpc_private_subnets_cidr_blocks" {
  value       = module.vpc.private_subnets_cidr_blocks
  description = "VPC Private Subnets CIDR Blocks"
}

output "vpc_database_subnets_cidr_blocks" {
  value       = module.vpc.database_subnets_cidr_blocks
  description = "VPC Database Subnets CIDR Blocks"
}

output "vpc_database_subnet_group_name" {
  value       = module.vpc.database_subnet_group_name
  description = "VPC Database Subnet Group Name"
}

output "public_subnet_id_1" {
  value       = data.aws_subnet.public[0].id
  description = "Public Subnet ID 1"
}

output "public_subnet_id_2" {
  value       = data.aws_subnet.public[1].id
  description = "Public Subnet ID 2"
}

output "private_subnet_id_1" {
  value       = data.aws_subnet.private[0].id
  description = "Private Subnet ID 1"
}

output "private_subnet_id_2" {
  value       = data.aws_subnet.private[1].id
  description = "Private Subnet ID 2"
}

