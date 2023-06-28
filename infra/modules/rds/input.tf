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

variable "subnet_ids" {
    type = list(string)
    description = "subnet ids"
}

# lambdaのセキュリティグループID
variable "lambda_security_group_id" {
  type        = string
  description = "lambda security group id"
}

# DBエンジン種類
variable "db_engine" {
  type        = string
  description = "DB engine type"
}
# DBエンジンバージョン
variable "db_engine_version" {
  type        = string
  description = "DB engine version"
}
# DBインスタンスタイプ
variable "db_instance_class" {
  type        = string
  description = "DB instance class"
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

variable "db_username" {
  type        = string
  description = "db user-name"
}

variable "db_name" {
  type        = string
  description = "(optional) describe your variable"
}




