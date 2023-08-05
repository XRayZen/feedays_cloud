package main

import (
	// "batch/ApiFunction"
	"batch/ApiFunction"
	"context"
	"log"

	// "os"

	// "encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	// "github.com/mitchellh/mapstructure"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// バッチ処理なのでリクエストをパースせずコードを実行するだけ
	dbRepo,inertVal,err := ApiFunction.InitDataBase(false)
	result, err := ApiFunction.Batch(dbRepo, inertVal)
	if err != nil || !result {
		log.Println("BATCH ERROR! :", err)
		return events.APIGatewayProxyResponse{
			Body:       string("BATCH ERROR! :" + err.Error()),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		Body:       string("BATCH SUCCESS!"),
		StatusCode: http.StatusOK,
	}, nil
}
