package test_lambda_function

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"test/api_gen_code"
)

// "encoding/json"
// "test/Data"

// "github.com/aws/aws-lambda-go/events"
// "github.com/aws/aws-lambda-go/lambda"

func execTest() (bool, error) {
	// Userのテスト
	result,userCfg,err := TestApiUser()
	if err != nil || !result {
		return false, err
	}
	// Siteのテスト
	result, sites, err := TestApiSiteSearch(userCfg.UserUniqueID)
	if err != nil {
		return false, err
	}
	if !result {
		log.Println("TestApiSiteSearch: Failed")
	}
	// 購読する
	result, err = TestApiSiteSubscribe(sites[0])
	if err != nil {
		return false, err
	}
	if !result {
		log.Println("TestApiSiteSubscribe: Failed")
	}
	result, err = TestApiSiteFetchArticles(sites[0])
	if err != nil {
		return false, err
	}
	if !result {
		log.Println("TestApiSiteFetchArticles: Failed")
	}
	return true, nil
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
