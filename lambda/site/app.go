package main

import (
	"context"
	"encoding/json"
	"net/http"

	"site/Repo"
	"site/RequestHandler"
	"site/api_gen_code"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	// "github.com/mitchellh/mapstructure"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// アクセスしてきたIPアドレスを取得する
	access_ip := request.RequestContext.Identity.SourceIP
	api_req, err := parseBodyRequest(request)
	if err != nil {
		return errorResponse(err, *api_req.RequestType, *api_req.UserId)
	}
	// 変換されたらリクエストタイプに応じて処理を分岐する
	// ここでDBに接続する
	db_repo := Repo.DBRepoImpl{}
	if err := db_repo.ConnectDB(false); err != nil {
		return errorResponse(err, *api_req.RequestType, *api_req.UserId)
	}
	res, err := RequestHandler.ParseRequestType(access_ip, db_repo, *api_req.RequestType, *api_req.UserId,
		*api_req.RequestArgumentJson1, *api_req.RequestArgumentJson2)
	if err != nil {
		return errorResponse(err, *api_req.RequestType, *api_req.UserId)
	}
	// レスポンスを返す
	return genApiResponse(*api_req.RequestType, *api_req.UserId, res), nil
}

// リクエストをボディからパース
func parseBodyRequest(request events.APIGatewayProxyRequest) (api_gen_code.PostSiteJSONRequestBody, error) {
	var api_req api_gen_code.PostSiteJSONRequestBody
	if err := json.Unmarshal([]byte(request.Body), &api_req); err != nil {
		return api_req, err
	}
	return api_req, nil
}

func genApiResponse(response_type string, user_id string, value string) events.APIGatewayProxyResponse {
	res := api_gen_code.APIResponse{
		RequestType:   &response_type,
		UserId:        &user_id,
		ResponseValue: &value,
	}
	res_str, err := json.Marshal(res)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Generate Response Error",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return events.APIGatewayProxyResponse{
		Body:       string(res_str),
		StatusCode: http.StatusOK,
	}
}

// エラーレスポンスを返す
func errorResponse(er error, RequestType string, userId string) (events.APIGatewayProxyResponse, error) {
	error_message := er.Error()
	response := api_gen_code.APIResponse{
		ResponseValue: &error_message,
		RequestType:   &RequestType,
		UserId:        &userId,
	}
	body, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Generate Response Error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: http.StatusInternalServerError,
	}, er
}
