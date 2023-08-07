package main

import (
	"context"
	"encoding/json"
	"net/http"

	// "write/DBRepo"
	// "write/RequestHandler"
	// "write/api_gen_code"
	"site/RequestHandler"
	"site/api_gen_code"
	"site/Repo"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mitchellh/mapstructure"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// アクセスしてきたIPアドレスを取得する
	access_ip := request.RequestContext.Identity.SourceIP
	var api_req api_gen_code.PostSiteJSONBody
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
	// ここでDIする
	res, err := RequestHandler.ParseRequestType(access_ip, Repo.DBRepoImpl{}, *api_req.RequestType, *api_req.UserId,
		*api_req.RequestArgumentJson1, *api_req.RequestArgumentJson2)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// ここでレスポンスを作る
	response := api_gen_code.APIResponse{
		RequestType:   api_req.RequestType,
		UserId:        api_req.UserId,
		ResponseValue: &res,
	}
	body, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
