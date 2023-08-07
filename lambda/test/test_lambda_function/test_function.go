package test_lambda_function

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"test/api_gen_code"
)

// "encoding/json"
// "test/Data"

// "github.com/aws/aws-lambda-go/events"
// "github.com/aws/aws-lambda-go/lambda"

func execTest() {
	// Userのテスト

	// Siteのテスト
	TestApiSiteSearch()

}

// リクエストを送信
func SendApiRequest(requestPostJson string, path string) (api_gen_code.APIResponse, error) {
	// APIエンドポイントにリクエストを送る
	endPoint := os.Getenv("API_ENDPOINT")
	// そこに/siteなどをつける
	endPoint = endPoint + "/" + path
	// HTTPポストリクエストを送る
	res, err := http.Post(endPoint, "application/json", strings.NewReader(requestPostJson))
	if err != nil {
		return api_gen_code.APIResponse{}, err
	}
	// レスポンスをパースする
	var response = api_gen_code.APIResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return api_gen_code.APIResponse{}, err
	}
	return response, nil
}
