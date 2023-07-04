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

output "vpc_database_subnets" {
  value       = module.vpc.database_subnets
  description = "VPC Database Subnets IDs"
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
