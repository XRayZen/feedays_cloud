# VPC関連
variable "vpc_id" {
  type        = string
  description = "(optional) describe your variable"
}
# 配置するavailability zone
variable "availability_zone" {
  type        = string
  description = "availability zone"
}

variable "vpc_cidr_block" {
  type        = string
  description = "(optional) describe your variable"
}

variable "vpc_private_subnets" {
  type        = list(string)
  description = "vpc private subnets"
}

variable "vpc_database_subnets" {
  type        = list(string)
  description = "vpc database subnets"
}

variable "vpc_private_subnets_cidr_blocks" {
  type        = list(string)
  description = "vpc private subnets cidr blocks"
}

variable "vpc_database_subnet_group_name" {
  type        = string
  description = "vpc database subnet group name"
}

variable "vpc_database_subnets_cidr_blocks" {
  type        = list(string)
  description = "vpc database subnets cidr blocks"
}

# DBエンジン種類
variable "db_engine" {
  type        = string
  description = "DB engine type"
  default     = "mysql"
}

# DBエンジンバージョン
variable "db_engine_version" {
  type        = string
  description = "DB engine version"
  default     = "8.0.21"
}

# DBインスタンスタイプ
variable "db_instance_class" {
  type        = string
  description = "DB instance class"
  default     = "db.t3.micro"
}

# DBパラメーターグループ（ファミリー）
variable "db_parameter_group_family" {
  type        = string
  description = "DB parameter group family"
  default     = "mysql8.0"
}
# DBストレージタイプ
variable "db_storage_type" {
  type        = string
  description = "DB storage type"
}
# DB割り当てストレージサイズ
variable "db_allocated_storage" {
  type        = string
  description = "DB allocated storage"
}
# DB最大割り当てストレージサイズ
variable "db_max_allocated_storage" {
  type        = string
  description = "DB max allocated storage"
}

# DBユーザー名
variable "db_username" {
  type        = string
  description = "db user-name"
}
# DB名
variable "db_name" {
  type        = string
  description = "(optional) describe your variable"
}

# RDS Proxyの設定
variable "rds_proxy_debug_logging" {
  type        = bool
  description = "rds proxy debug logging"
  default     = false
}

variable "rds_proxy_engine_family" {
  type        = string
  description = "rds proxy engine family"
  default     = "MYSQL"
}

# RDSインスタンスの設定
variable "rds_create_cloudwatch_log_group" {
  type        = bool
  description = "rds create cloudwatch log group"
  default     = true
}

variable "rds_enabled_cloudwatch_logs_exports" {
  type        = list(string)
  description = "rds enabled cloudwatch logs exports"
  default     = ["audit", "error", "general", "slowquery"]
}

variable "tags" {
  type        = map(string)
  description = "tags"
}

variable "secret_version_stages" {
  type        = set(string)
  description = "database secret version stages"
}

