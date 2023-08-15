package main

import (
	"context"
	"encoding/json"
	"net/http"
	"read/Repo"
	"read/RequestHandler"
	"read/api_gen_code"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	// "github.com/mitchellh/mapstructure"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	api_req,err := parseBodyRequest(request)
	if err != nil {
		return errorResponse(err, "ParseRequest", "Parse Request Error")
	}
	// 変換されたらリクエストタイプに応じて処理を分岐する
	db_repo := Repo.DBRepoImpl{}
	if err := db_repo.ConnectDB(false); err != nil {
		return errorResponse(err, "ConnectDB", *api_req.UserId)
	}
	res, err := RequestHandler.ParseRequestType(db_repo, *api_req.RequestType, *api_req.UserId)
	if err != nil {
		return errorResponse(err, "ParseRequestType", *api_req.UserId)
	}
	// ここでレスポンスを作る
	return genApiResponse(*api_req.RequestType, *api_req.UserId, res), nil
}

// リクエストをボディからパース
func parseBodyRequest(request events.APIGatewayProxyRequest) (api_gen_code.PostReadJSONRequestBody, error) {
	var api_req api_gen_code.PostReadJSONRequestBody
	if err := json.Unmarshal([]byte(request.Body), &api_req); err != nil {
		return api_gen_code.PostReadJSONRequestBody{}, err
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
