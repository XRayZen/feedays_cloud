
# APIGWのエンドポイント
output "api_gw_endpoint" {
    value = "${aws_api_gateway_deployment.api_gw_deployment.invoke_url}"
    description = "API Gateway Endpoint"
}

