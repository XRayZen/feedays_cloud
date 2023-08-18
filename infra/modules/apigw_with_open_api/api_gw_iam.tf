
# Lambda/read関数に対するInvoke権限を与える
resource "aws_lambda_permission" "apigw_read_permission" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_read_arn
  # 許可を与えるAWSサービス
  principal = "apigateway.amazonaws.com"

  # API GatewayのREST APIのIDを指定する
  # このIDは、API GatewayのREST APIを作成するときに出力される
  # または、API Gatewayのコンソール画面から確認できる
  source_arn = "${aws_api_gateway_rest_api.api_gw.execution_arn}/*"
}

# Lambda/user関数に対するInvoke権限を与える
resource "aws_lambda_permission" "apigw_user_permission" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_user_arn
  # 許可を与えるAWSサービス
  principal = "apigateway.amazonaws.com"

  # API GatewayのREST APIのIDを指定する
  # このIDは、API GatewayのREST APIを作成するときに出力される
  # または、API Gatewayのコンソール画面から確認できる
  source_arn = "${aws_api_gateway_rest_api.api_gw.execution_arn}/*"
}

# Lambda/site関数に対するInvoke権限を与える
resource "aws_lambda_permission" "apigw_site_permission" {
    statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_site_arn
  # 許可を与えるAWSサービス
  principal = "apigateway.amazonaws.com"

  # API GatewayのREST APIのIDを指定する
  # このIDは、API GatewayのREST APIを作成するときに出力される
  # または、API Gatewayのコンソール画面から確認できる
  source_arn = "${aws_api_gateway_rest_api.api_gw.execution_arn}/*"
}

################################
# API GatewayにアタッチするIAM Role
################################
data "aws_iam_policy_document" "api_gw_assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["apigateway.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "api_gateway_role" {
  name               = "${var.api_gw_name}-apigateway-role"
  assume_role_policy = data.aws_iam_policy_document.api_gw_assume_role.json
}

resource "aws_iam_role_policy_attachment" "api_gateway_policy_logs" {
  role       = aws_iam_role.api_gateway_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs"
}

resource "aws_iam_role_policy_attachment" "api_gateway_policy_lambda" {
  role       = aws_iam_role.api_gateway_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaRole"
}
