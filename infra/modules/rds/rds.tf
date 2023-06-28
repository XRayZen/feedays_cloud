
# DBを作成
resource "aws_db_instance" "mysql_db_instance" {
  # DBの基本設定
  name                 = var.db_name
  engine               = var.db_engine
  engine_version       = var.db_engine_version
  instance_class       = var.db_instance_class
  username             = var.db_username
  password             = random_password.db-password.result
  parameter_group_name = "default.mysql8.0"

  # ストレージ設定
  allocated_storage     = var.db_allocated_storage
  max_allocated_storage = var.db_max_allocated_storage
  storage_type          = "gp2"

  # ネットワーク設定
  vpc_security_group_ids = ["${aws_security_group.mysql-sg.id}"]
  db_subnet_group_name   = aws_db_subnet_group.mysql_db_subnet_group.name
  # publicアクセスを許可するか設定
  publicly_accessible = false
  # ポート番号を設定
  port = 3306

  # マルチAZしないのならfalseにして配置するavailability zoneを指定
  multi_az          = false
  availability_zone = var.availability_zone

  # -- DBの管理設定。maintenance_windowは、backup_windowの後の時間に設定する
  # バックアップを行う時間を設定
  backup_window = "04:00-05:00"
  # バックアップの保存期間（日）を設定
  backup_retention_period = 0
  # DBインスタンスまたはクラスターのエンジンバージョンの更新、OS更新があった場合に更新作業を行う時間を設定
  maintenance_window = "Mon:05:00-Mon:08:00"
  # 自動的にDBのマイナーバージョンアップグレードを行うか設定する
  auto_minor_version_upgrade = false

  # -- 削除設定
  # 削除操作を受付るかを指定。削除させない場合はtrue
  deletion_protection = false
  # インスタンス削除時にスナップショットをとるかを設定
  skip_final_snapshot = true
  # DBインスタンスが削除されたときに保存するスナップショットの名前 skip_final_snapshot = falseの時に指定
  final_snapshot_identifier = "final-snapshot-lab"
  # データベースの変更をすぐに適用するか、次のメンテナンスウィンドウ中に適用するかを指定する
  apply_immediately = true
  tags = {
    Name = "${var.db_name}"
  }
}

# RDS DBサブネットグループを作成
resource "aws_db_subnet_group" "mysql_db_subnet_group" {
  name = "${var.db_name}-subnet-group"
  subnet_ids = [
    "${var.private_subnet_id_1}",
    "${var.private_subnet_id_2}"
  ]
  tags = {
    Name = "${var.db_name}-subnet-group"
  }
}

# RDS DBセキュリティグループを作成
resource "aws_security_group" "mysql-sg" {
  name        = "${var.db_name}-security-group"
  description = "MySQL security group"
  vpc_id      = var.vpc_id

  ingress {
    description = "MySQL from VPC"
    from_port   = 3306
    to_port     = 3306
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr}"]
  }

  ingress {
    description = "MySQL from rds Proxy"
    from_port   = 3306
    to_port     = 3306
    protocol    = "tcp"
    security_groups = ["${aws_security_group.rds_proxy.id}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.db_name}-security-group"
  }
}














