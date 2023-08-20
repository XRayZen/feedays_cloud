
resource "aws_api_gateway_rest_api" "api_gw" {
  name = var.api_gw_name
  body = templatefile("./OpenAPI/backend_api.yaml",
    {
      lambda_read_arn     = var.lambda_read_arn,
      lambda_user_arn     = var.lambda_user_arn,
      lambda_site_arn     = var.lambda_site_arn,
      apigateway_role_arn = "${aws_iam_role.api_gateway_role.arn}"
  })
  description = var.api_gw_description
  put_rest_api_mode = "merge"

  endpoint_configuration {
    types = var.endpoint_configuration_types
  }
  # policy属性が変更された場合に、リソースを再作成しないように指定
  lifecycle {
    ignore_changes = [
      policy
    ]
  }
}

resource "aws_api_gateway_deployment" "api_gw_deployment" {
  depends_on = [
    aws_api_gateway_rest_api.api_gw
  ]

  rest_api_id = aws_api_gateway_rest_api.api_gw.id
  stage_name  = var.stage_name

  # lifecycle {
  #   create_before_destroy = true
  # }

  triggers = {
    redeployment = "${sha1(file("./OpenAPI/backend_api.yaml"))}"
  }
}

resource "aws_api_gateway_usage_plan" "api_usage_plan" {
  name        = "${var.project_name}_rest_api_usage_plan"
  description = "api_usage_plan"
  api_stages {
    api_id = aws_api_gateway_rest_api.api_gw.id
    stage  = aws_api_gateway_deployment.api_gw_deployment.stage_name
  }
  product_code = "${var.project_name}_rest_api"

  quota_settings {
    # 時間内に行うことができる最大リクエスト数
    limit = var.max_request_limit
    # 初期時間帯に与えられた制限から減算されるリクエスト数
    offset = 0
    # 制限が適用される時間帯。有効な値は"DAY"、 "WEEK"、または"MONTH"
    period = var.limit_period
  }
  throttle_settings {
    # APIリクエストバースト制限
    burst_limit = var.burst_limit
    # APIリクエストのレート制限
    rate_limit = var.rate_limit
  }
}
