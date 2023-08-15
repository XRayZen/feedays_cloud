package main

import (
	"context"
	"fmt"
	"log"
	FetchSecret "test/FetchSecret"
	Internet "test/internet"
	RDS "test/rds"
	"test/test_lambda_function"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(HandleRequest)
}

func LogWrite(msg string) {
	if RDS.Debug {
		log.Println(msg)
	}
}

func HandleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	// リクエストをパースせずにテストを実行する
	// APIが要求する機能を全て仮実装してテスト検証をする
	log.Default()
	log.Println("Secret Read Test Start")
	result, err := FetchSecret.Secret_read_test()
	if err != nil || !result {
		return genLambdaResponse("Failed", "Secret Read Test Failed :"+err.Error()),nil
	}
	fmt.Println("Secret Read Test Success")
	// データベース読み書きテスト
	result, err = RDS.RdsWriteReadTest()
	if err != nil || !result {
		return genLambdaResponse("Failed", "RDS Write Read Test Failed :"+err.Error()),nil
	}
	fmt.Println("RDS Write Read Test Success")
	fmt.Println("Internet Connection Test Start")
	// インターネット導通テスト
	str, err := Internet.GetGIGAZINE()
	if err != nil || str == "" {
		return genLambdaResponse("Failed", "Internet Connection Test Failed err : "+err.Error()),nil
	}
	log.Println("Gigazine RSS Title: " + str)
	fmt.Println("Internet Connection Test Success")
	// LambdaAPIのテスト
	log.Println("Lambda API Test Start")
	if result ,err := test_lambda_function.LambdaApiTest(); err != nil || !result {
		// テスト失敗
		return genLambdaResponse("Failed", "Lambda API Test Failed err : "+err.Error()),nil
	}
	// 全てのテストが成功したら200を返す
	return genLambdaResponse("Success", "All Test Success"), nil
}

func genLambdaResponse(response_type string, message string) events.LambdaFunctionURLResponse {
	switch response_type {
	case "Success":
		return events.LambdaFunctionURLResponse{
			StatusCode:      200,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            message,
			IsBase64Encoded: false,
			Cookies: []string{
				"cookie1",
			},
		}
	case "Failed":
		return events.LambdaFunctionURLResponse{
			StatusCode:      400,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            message,
			IsBase64Encoded: false,
			Cookies: []string{
				"cookie1",
			},
		}
	}
	return events.LambdaFunctionURLResponse{
		StatusCode:      400,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            "error msg : " + "response_type is not Success or Failed",
		IsBase64Encoded: false,
		Cookies: []string{
			"cookie1",
		},
	}
}
