package main

import (
	"context"
	"encoding/json"
	"net/http"

	"site/Data"
	"site/Repo"
	"site/RequestHandler"
	"site/api_gen_code"

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
	// mapstructureではなくBodyをデコードした方がいいのではないか
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
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}
	// レスポンスを返す
	response := api_gen_code.APIResponse{
		RequestType:   api_req.RequestType,
		UserId:        api_req.UserId,
		ResponseValue: &res,
	}
	res_json, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       string(res_json),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}

func GenAPIResponse(responseType string, value string, errorMsg string) (string, error) {
	res := Data.APIResponse{
		ResponseType: responseType,
		Value:        value,
		Error:        errorMsg,
	}
	res_str, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(res_str), nil
}

