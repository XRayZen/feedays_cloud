
output "rds_proxy_endpoint" {
  value = module.rds_proxy.proxy_endpoint
}

# 試しに作ってみる
output "rds_proxy_read_write_endpoint" {
  value = module.rds_proxy.db_proxy_endpoints["read_write"].endpoint
}

output "rds_proxy_arn"{
  value = module.rds_proxy.proxy_arn
}

output "db_password" {
  value = random_password.db-password.result
  sensitive = true
}
