package main

import (
	"context"
	Internet "test/internet"
	RDS "test/rds"
	SecretManager "test/secret_manager"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	// リクエストをパースせずにテストを実行する
	// シークレットリードテスト
	result, err := SecretManager.Secret_read_test()
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode:      400,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            "error msg: " + err.Error(),
			IsBase64Encoded: false,
			Cookies: []string{
				"cookie1",
			},
		}, err
	}
	if !result {
		return events.LambdaFunctionURLResponse{
			StatusCode:      400,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            "Secret Read Test Failed",
			IsBase64Encoded: false,
			Cookies: []string{
				"cookie1",
			},
		}, nil
	}
	// データベース読み書きテスト
	result, err = RDS.RdsWriteReadTest()
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode:      400,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            "error msg: " + err.Error(),
			IsBase64Encoded: false,
			Cookies: []string{
				"cookie1",
			},
		}, err
	}
	if !result {
		return events.LambdaFunctionURLResponse{
			StatusCode:      400,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            "RDS Write Read Test Failed",
			IsBase64Encoded: false,
			Cookies: []string{
				"cookie1",
			},
		}, nil
	}
	// インターネット導通テスト
	str, err := Internet.GetGIGAZINE()
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode:      400,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            "error msg: " + err.Error(),
			IsBase64Encoded: false,
			Cookies: []string{
				"cookie1",
			},
		}, err
	}
	if str == "" {
		return events.LambdaFunctionURLResponse{
			StatusCode:      400,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            "Internet Connection Test Failed",
			IsBase64Encoded: false,
			Cookies: []string{
				"cookie1",
			},
		}, nil
	}
	// 全てのテストが成功したら200を返す
	return events.LambdaFunctionURLResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            "All Test Success",
		IsBase64Encoded: false,
		Cookies: []string{
			"cookie1",
		},
	}, nil
}
