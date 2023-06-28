#--------------------------------------------------------------
# Secrets Manager
#--------------------------------------------------------------

resource "random_password" "db-password" {
  length           = 16
  special          = true
  override_special = "_!%^"
}

resource "aws_secretsmanager_secret" "db-password" {
  name = "${var.name}-db-password"
}

resource "aws_secretsmanager_secret_version" "db-password" {
  secret_id     = aws_secretsmanager_secret.db-password.id
  secret_string = random_password.db-password.result
}

data "aws_secretsmanager_secret" "db-password" {
  name       = "${var.name}-db-password"
  depends_on = [aws_secretsmanager_secret.db-password]
}

data "aws_secretsmanager_secret_version" "db-password" {
  secret_id  = data.aws_secretsmanager_secret.db-password.id
  depends_on = [aws_secretsmanager_secret_version.db-password]
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

