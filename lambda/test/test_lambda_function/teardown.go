package test_lambda_function

import (
	"encoding/json"
	"log"
	"test/api_gen_code"
)

// テストで使われたデータを全て物理削除モードで削除する
func DeleteTestData(user_id string , site_url string) (bool, error) {
	request_type := "DeleteUserData"
	is_unscoped_json, _ := json.Marshal(true)
	is_unscoped_json_str := string(is_unscoped_json)
	post_request := api_gen_code.PostSiteJSONRequestBody{
		RequestType:          &request_type,
		UserId:               &user_id,
		RequestArgumentJson1: &is_unscoped_json_str,
	}
	post_request_json, _ := json.Marshal(post_request)
	response, err := SendApiRequest(string(post_request_json), "user")
	if err != nil ||*response.ResponseValue != "Success DeleteUserData"  {
		log.Println("Failed DeleteUserData : ", err)
		return false, err
	}
	log.Println("Success DeleteUserData")
	// Siteを削除する
	request_type = "DeleteSite"
	is_unscoped_json, _ = json.Marshal(true)
	is_unscoped_json_str = string(is_unscoped_json)
	site_url_json, _ := json.Marshal(site_url)
	site_url_json_str := string(site_url_json)
	post_request = api_gen_code.PostSiteJSONRequestBody{
		RequestType:          &request_type,
		UserId:               &user_id,
		RequestArgumentJson2: &site_url_json_str,
		RequestArgumentJson1: &is_unscoped_json_str,
	}
	post_request_json, _ = json.Marshal(post_request)
	response, err = SendApiRequest(string(post_request_json), "site")
	if err != nil ||*response.ResponseValue != "Success DeleteSite"  {
		log.Println("Failed DeleteSite : ", err)
		return false, err
	}
	log.Println("Success DeleteSite")
	return true, nil
}
