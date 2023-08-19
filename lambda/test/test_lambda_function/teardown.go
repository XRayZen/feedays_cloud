package test_lambda_function

import (
	"encoding/json"
	"log"
	"test/api_gen_code"
)

// サービス終了処理で作ったテストデータをテーブルごと削除する
func DeleteTestData(user_id string, site_url string) (bool, error) {
	request_type := "ServiceFinalize"
	request := api_gen_code.PostSiteJSONRequestBody{
		RequestType: &request_type,
	}
	request_json, _ := json.Marshal(request)
	result, err := SendApiRequest(string(request_json), "site")
	if err != nil {
		log.Println("Failed Delete TestData")
		return false, err
	}
	expected := "Success ServiceFinalize"
	if result.ResponseValue != &expected {
		log.Println("Failed Delete TestData")
		return false, nil
	}
	log.Println("Success Delete TestData")
	return true, nil
}

