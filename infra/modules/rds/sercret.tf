#--------------------------------------------------------------
# Secrets Manager
#--------------------------------------------------------------
# data "aws_kms_alias" "secretsmanager" {
#   name = "alias/aws/secretsmanager"
# }

locals {
  db_password = random_password.db-password.result
}

resource "random_password" "db-password" {
  length           = 16
  special          = true
  override_special = "_!%^"
}

resource "aws_secretsmanager_secret" "superuser" {
  name        = var.db_username
  description = "Database superuser, ${var.db_username}, database connection values"
  # kms_key_id  = data.aws_kms_alias.secretsmanager.id

  tags = var.tags
}

resource "aws_secretsmanager_secret_version" "superuser" {
  secret_id = aws_secretsmanager_secret.superuser.id
  version_stages = var.secret_version_stages
  secret_string = jsonencode({
    username = var.db_username
    password = local.db_password
  })
}

# これでダメだったらjson形式にしておく
# The map here can come from other supported configurations
# like locals, resource attribute, map() built-in, etc.
# variable "example" {
#   default = {
#     key1 = "value1"
#     key2 = "value2"
#   }

#   type = map(string)
# }

# resource "aws_secretsmanager_secret_version" "example" {
#   secret_id     = aws_secretsmanager_secret.example.id
#   secret_string = jsonencode(var.example)
# }

