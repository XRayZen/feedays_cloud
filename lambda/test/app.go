package main

import (
	"context"
	"fmt"
	"log"
	FetchSecret "test/fetch_secret"
	Internet "test/internet"
	RDS "test/rds"

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
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode:      400,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            "error msg : " + err.Error(),
			IsBase64Encoded: false,
			Cookies: []string{
				"cookie1",
			},
		}, err
	}
	log.Println("Secret Read Test End")
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
	fmt.Println("Secret Read Test Success")
	// データベース読み書きテスト
	result, err = RDS.RdsWriteReadTest()
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode:      400,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            "error msg : " + err.Error(),
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
	fmt.Println("RDS Write Read Test Success")
	fmt.Println("Internet Connection Test Start")
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
	fmt.Println("Internet Connection Test End")
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
	log.Println("Gigazine RSS Title: " + str)
	fmt.Println("Internet Connection Test Success")
	// 入れ子の検索テスト
	result, err = RDS.DbNestedStructTest()
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode:      400,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            "error msg : " + err.Error(),
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
			Body:            "RDS Nested Struct Test Failed",
			IsBase64Encoded: false,
			Cookies: []string{
				"cookie1",
			},
		}, nil
	}
	fmt.Println("RDS Nested Struct Test Success")
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