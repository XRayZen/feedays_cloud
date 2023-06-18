package main

import (
	"context"
	// "encoding/json"
	"net/http"

	// "read/Repo"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	// "github.com/mitchellh/mapstructure"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// バッチ処理なのでリクエストを受け付けずコードを実行するだけ
	
	return events.APIGatewayProxyResponse{
		Body:       string(""),
		StatusCode: http.StatusOK,
	}, nil
}


func main() {
	lambda.Start(HandleRequest)
}
