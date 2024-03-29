
variable "project_name" {
  type        = string
  description = "Project Name"
}

variable "api_gw_name" {
  type        = string
  description = "API Gateway Name"
}

variable "api_gw_description" {
  type        = string
  description = "API Gateway Description"
}

variable "endpoint_configuration_types" {
  type        = list(string)
  description = "Endpoint Configuration Types"
  default     = ["REGIONAL"]
}

variable "stage_name" {
  type        = string
  description = "Stage Name"
}

# 時間内に行うことができる最大リクエスト数
# limit = 1000
variable "max_request_limit" {
  type        = number
  description = "Max Request Limit"
  default     = 1000
}

# 制限が適用される時間帯。有効な値は"DAY"、 "WEEK"、または"MONTH"
variable "limit_period" {
  type        = string
  description = "Limit Period"
  default     = "DAY"
}

# APIリクエストバースト制限
# burst_limit = 100
variable "burst_limit" {
  type        = number
  description = "Burst Limit"
  default     = 100
}
# APIリクエストのレート制限
# rate_limit = 50
variable "rate_limit" {
  type        = number
  description = "Rate Limit"
  default     = 50
}

# Lambda関数のARN
variable "lambda_read_arn" {
  type        = string
  description = "Lambda ARN"
}

variable "lambda_user_arn" {
  type        = string
  description = "Lambda ARN"
}

variable "lambda_site_arn" {
  type        = string
  description = "Lambda ARN"
}

variable "lambda_read_name" {
  type        = string
  description = "Lambda Read Name"
}

variable "lambda_site_name" {
  type        = string
  description = "Lambda Site Name"
}

variable "lambda_user_name" {
  type        = string
  description = "Lambda User Name"
}

