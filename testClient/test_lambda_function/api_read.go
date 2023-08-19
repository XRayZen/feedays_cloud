package test_lambda_function

import (
	"encoding/json"
	"log"
	"user/Data"
	"user/api_gen_code"
)

func TestApiRead(userId string) (bool, Data.UserConfig, error) {
	result , err := testApiReadExploreCategories(userId)
	if err != nil || !result {
		log.Println("TestApiReadExploreCategories: Failed")
		return false, Data.UserConfig{}, err
	}
	return true, Data.UserConfig{}, nil
}

// ExploreCategories
func testApiReadExploreCategories(userId string) (bool, error) {
	req_type := "ExploreCategories"
	result, err := SendApiReadRequest(api_gen_code.PostReadJSONRequestBody{
		RequestType: &req_type,
		UserId:      &userId,
	})
	if err != nil {
		return false, err
	}
	// 結果をパースする
	var categories []Data.ExploreCategory
	err = json.Unmarshal([]byte(result), &categories)
	if err != nil {
		return false, err
	}
	if len(categories) == 0 {
		return false, nil
	}
	return true, nil
}

func SendApiReadRequest(request api_gen_code.PostReadJSONRequestBody) (string, error) {
	json, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	result, err := SendApiRequest(string(json), "/read")
	if err != nil {
		return "", err
	}
	return *result.ResponseValue, nil
}
