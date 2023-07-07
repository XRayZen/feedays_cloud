# 秘密鍵のアルゴリズム設定
resource "tls_private_key" "keygen" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

module "ec2_key_pair" {
  source = "terraform-aws-modules/key-pair/aws"

  key_name   = var.key_name
  public_key = tls_private_key.keygen.public_key_openssh
}

# クライアントPCにKey pair（秘密鍵と公開鍵）を作成
locals {
  # クライアントPCの公開鍵はstage/develop/test_ec2に保存される
  public_key_file = "../../../${var.key_name}.id_rsa.pub"
  # クライアントPCの秘密鍵はstage/develop/test_ec2に保存される
  private_key_file = "../../../${var.key_name}.pem"
}

#local_fileのリソースを指定するとterraformを実行するディレクトリ内でファイル作成やコマンド実行が出来る。
resource "local_file" "public_key_openssh" {
  content  = tls_private_key.keygen.public_key_openssh
  filename = local.public_key_file
  # ファイルのパーミッションを600に変更
  provisioner "local-exec" {
    command = "chmod 600 ${local.public_key_file}"
  }
}

resource "local_file" "private_key_pem" {
  content  = tls_private_key.keygen.private_key_pem
  filename = local.private_key_file
  # ファイルのパーミッションを600に変更
  provisioner "local-exec" {
    command = "chmod 600 ${local.private_key_file}"
  }
}

