package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"user/DbRepo"
	"user/RequestHandler"
	"user/api_gen_code"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// アクセスしてきたIPアドレスを取得する
	access_ip := request.RequestContext.Identity.SourceIP
	log.Println("Access IP: ", access_ip)
	log.Println("Request Body: ", request.Body)
	// リクエストを変換する為にPostUserJSONBodyを使う
	api_req, err := parseBodyRequest(request)
	if err != nil {
		return errorResponse(err, *api_req.RequestType, *api_req.UserId)
	}
	// 変換されたらリクエストタイプに応じて処理を分岐する
	db_repo := DbRepo.DBRepoImpl{}
	if err := db_repo.ConnectDB(false); err != nil {
		return errorResponse(err, *api_req.RequestType, *api_req.UserId)
	}
	// エラーの原因はnullを参照したからだろう
	res, err := RequestHandler.ParseRequestType(access_ip, db_repo, *api_req.RequestType, *api_req.UserId,
		*api_req.RequestArgumentJson1, *api_req.RequestArgumentJson2)
	if err != nil {
		return errorResponse(err, *api_req.RequestType, *api_req.UserId)
	}
	// ここでレスポンスを作る
	return genApiResponse(*api_req.RequestType, *api_req.UserId, res), nil
}

// リクエストをボディからパース
func parseBodyRequest(request events.APIGatewayProxyRequest) (api_gen_code.PostUserJSONBody, error) {
	var api_req api_gen_code.PostUserJSONBody
	if err := json.Unmarshal([]byte(request.Body), &api_req); err != nil {
		return api_gen_code.PostUserJSONBody{}, err
	}
	// nullを参照するとエラーになるので、nullの場合は空文字にする
	if api_req.UserId == nil {
		api_req.UserId = new(string)
		*api_req.UserId = ""
	}
	if api_req.RequestArgumentJson1 == nil {
		api_req.RequestArgumentJson1 = new(string)
		*api_req.RequestArgumentJson1 = ""
	}
	if api_req.RequestArgumentJson2 == nil {
		api_req.RequestArgumentJson2 = new(string)
		*api_req.RequestArgumentJson2 = ""
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
