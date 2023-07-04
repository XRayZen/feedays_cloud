
output "rds_proxy_endpoint" {
  value = module.rds_proxy.proxy_endpoint
}

output "rds_proxy_arn"{
  value = module.rds_proxy.proxy_arn
}

output "db_password" {
  value = random_password.db-password.result
  sensitive = true
}
