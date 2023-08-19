include "root" {
    path= find_in_parent_folders()
}

locals{
    env= read_terragrunt_config(find_in_parent_folders("env.hcl"))
}

terraform{
    source = "../../../modules/apigw_with_open_api"
}

dependency "lambda_read"{
    config_path = "../lambda_read"

    mock_outputs = {
        lambda_function_invoke_arn = "arn:aws:lambda:us-east-1:123456789012:function:lambda_read"
        lambda_function_name = "lambda_read"
    }
}

dependency "lambda_site"{
    config_path = "../lambda_site"

    mock_outputs = {
        lambda_function_invoke_arn = "arn:aws:lambda:us-east-1:123456789012:function:lambda_site"
        lambda_function_name = "lambda_site"
    }
}

dependency "lambda_user"{
    config_path = "../lambda_user"

    mock_outputs = {
        lambda_function_invoke_arn = "arn:aws:lambda:us-east-1:123456789012:function:lambda_user"
        lambda_function_name = "lambda_user"
    }
}

inputs = {
    project_name = local.env.locals.project_name
    api_gw_name = "API_Gateway_for_${local.env.locals.project_name}_${local.env.locals.stage}"
    endpoint_configuration_types = ["REGIONAL"]
    stage_name = local.env.locals.stage
    api_gw_description = "API Gateway for ${local.env.locals.project_name} ${local.env.locals.stage}"
    # APIリクエスト制限
    max_request_limit = 1000
    limit_period = "DAY"
    burst_limit = 200
    rate_limit = 100
    # Lambda ARN
    lambda_read_arn = dependency.lambda_read.outputs.lambda_function_invoke_arn
    lambda_site_arn = dependency.lambda_site.outputs.lambda_function_invoke_arn
    lambda_user_arn = dependency.lambda_user.outputs.lambda_function_invoke_arn
    lambda_read_name = dependency.lambda_read.outputs.lambda_function_name
    lambda_site_name = dependency.lambda_site.outputs.lambda_function_name
    lambda_user_name = dependency.lambda_user.outputs.lambda_function_name
}
