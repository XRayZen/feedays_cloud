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
	// "github.com/mitchellh/mapstructure"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// アクセスしてきたIPアドレスを取得する
	access_ip := request.RequestContext.Identity.SourceIP
	api_req,err := decodeApiRequest(request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}
	// 変換されたらリクエストタイプに応じて処理を分岐する
	// ここでDBに接続する
	db_repo := Repo.DBRepoImpl{}
	if err := db_repo.ConnectDB(false); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}
	res, err := RequestHandler.ParseRequestType(access_ip, db_repo, *api_req.RequestType, *api_req.UserId,
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
		Headers: map[string]string{
			"application/json": "charset=utf-8",
		},
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

func decodeApiRequest(request events.APIGatewayProxyRequest) (api_gen_code.PostSiteJSONRequestBody, error) {
	var api_req api_gen_code.PostSiteJSONRequestBody
	if err := json.Unmarshal([]byte(request.Body), &api_req); err != nil {
		return api_gen_code.PostSiteJSONRequestBody{}, err
	}
	return api_req, nil
}
