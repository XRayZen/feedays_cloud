package main

import (
	"fmt"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"encoding/json"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/mitchellh/mapstructure"
	"read/api_gen_code"
)


func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var api_req api_gen_code.PostReadJSONBody
	decoderConfig := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &api_req,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	err = decoder.Decode(request.QueryStringParameters)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// 変換されたらリクエストタイプに応じて処理を分岐する
	// 別のパッケージに移して処理を書く
	
}

func main() {
	lambda.Start(HandleRequest)
}


