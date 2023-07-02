variable "lambda_function_name" {
  type        = string
  description = "lambda function name"
}

variable "lambda_function_description" {
  type        = string
  description = "lambda function description"
}

variable "lambda_function_architecture" {
  type        = string
  description = "lambda function architecture"
  default     = "x86_64"
}

variable "repo_url" {
  type        = string
  description = "ecr url"
}

variable "image_tag" {
  type = string
  # タグを指定する
  description = "image tag (e.g. latest)"
}

variable "memory_size" {
  type        = number
  description = "memory size"
}

variable "timeout" {
  type        = number
  description = "timeout"
}

# 環境変数はTerragruntで設定する
variable "variables" {
  type        = map(string)
  description = "environment variables"
}

variable "managed_policy_arns" {
  type        = set(string)
  description = "managed policy arns"
}

variable "subnet_ids" {
  type        = list(string)
  description = "subnet ids"
}

variable "vpc_id" {
  type        = string
  description = "vpc id"
}

variable "vpc_private_subnets_cidr_blocks" {
  type        = list(string)
  description = "vpc private subnets cidr blocks"
}

variable "vpc_database_subnets_cidr_blocks" {
  type        = list(string)
  description = "vpc database subnets cidr blocks"
}

