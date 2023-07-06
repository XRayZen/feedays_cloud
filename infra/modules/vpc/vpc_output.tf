output "vpc_id" {
  value       = module.vpc.vpc_id
  description = "VPC ID"
}

output "vpc_cidr" {
  value       = module.vpc.vpc_cidr_block
  description = "VPC CIDR"
}

output "vpc_public_subnets" {
  value = module.vpc.public_subnets
  description = "VPC Public Subnets IDs"
}

output "vpc_private_subnets" {
  value       = module.vpc.private_subnets
  description = "VPC Private Subnets IDs"
}

output "vpc_database_subnets" {
  value       = module.vpc.database_subnets
  description = "VPC Database Subnets IDs"
}

output "vpc_public_subnets_cidr_blocks" {
  value = module.vpc.public_subnets_cidr_blocks
  description = "VPC Public Subnets CIDR Blocks"
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

# 個別のサブネットID
output "vpc_public_subnet_1_id" {
  value       = data.aws_subnet.public[0].id
  description = "VPC Public Subnet 1 ID"
}
