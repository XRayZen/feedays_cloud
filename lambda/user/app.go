package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"user/DBRepo"
	"user/RequestHandler"
	"user/api_gen_code"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// アクセスしてきたIPアドレスを取得する
	access_ip := request.RequestContext.Identity.SourceIP
	// リクエストを変換する為にPostUserJSONBodyを使う
	var api_req api_gen_code.PostUserJSONBody
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
	// ここでDIする
	dbRepo := DBRepo.DBRepoImpl{}
	if err := dbRepo.ConnectDB(false); err != nil {
		return errorResponse(err, *api_req.RequestType, *api_req.UserId)
	}
	if err := dbRepo.AutoMigrate(); err != nil {
		return errorResponse(err, *api_req.RequestType, *api_req.UserId)
	}
	res, err := RequestHandler.ParseRequestType(access_ip, dbRepo, *api_req.RequestType, *api_req.UserId,
		*api_req.RequestArgumentJson1, *api_req.RequestArgumentJson2)
	if err != nil {
		return errorResponse(err, *api_req.RequestType, *api_req.UserId)
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

// エラーレスポンスを返す
func errorResponse(er error, RequestType string, userId string) (events.APIGatewayProxyResponse, error) {
	errorMessage := er.Error()
	response := api_gen_code.APIResponse{
		ResponseValue: &errorMessage,
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
