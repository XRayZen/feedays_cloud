package test_lambda_function

import (
	"encoding/json"
	"log"
	"user/api_gen_code"
)

// サービス終了処理で作ったテストデータをテーブルごと削除する
func DeleteTestData(user_id string, site_url string) (bool, error) {
	request_type := "ServiceFinalize"
	request := api_gen_code.PostUserJSONRequestBody{
		RequestType: &request_type,
	}
	request_json, _ := json.Marshal(request)
	result, err := SendApiRequest(string(request_json), "user")
	if err != nil {
		log.Println("Failed Delete TestData")
		return false, err
	}
	answer := result.ResponseValue
	an := "" + *answer
	expected := "Success ServiceFinalize"
	if an != expected {
		log.Println("Failed Delete TestData")
		return false, nil
	}
	log.Println("Success Delete TestData")
	return true, nil
}
