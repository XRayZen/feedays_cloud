package test_lambda_function

import (
	"encoding/json"
	"log"

	"test/Data"
	"test/api_gen_code"
	"time"
)

func TestApiUserPart1() (bool, Data.UserConfig, error) {
	request_type := "GenUserID"
	request_response, err := SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType: &request_type,
	})
	if err != nil || len(request_response) < 10 {
		log.Println("Failed to get user id")
		return false, Data.UserConfig{}, err
	}
	log.Println("Success to get user id")
	user_Id := request_response
	// ユーザー登録
	result, err := testRegisterUser(user_Id)
	if err != nil || !result {
		log.Println("Failed to register user")
		return false, Data.UserConfig{}, err
	}
	// ConfigSync
	result_config_sync, err := testConfigSync(user_Id)
	if err != nil || !result_config_sync {
		log.Println("Failed to config sync")
		return false, Data.UserConfig{}, err
	}
	// ReportReadActivity
	result, err = testReportReadActivity(user_Id)
	if err != nil || !result {
		log.Println("Failed to report read activity")
		return false, Data.UserConfig{}, err
	}
	// UpdateConfig
	result, user_config, err := testUpdateConfig(user_Id)
	if err != nil || !result {
		log.Println("Failed to update config")
		return false, Data.UserConfig{}, err
	}
	// ModifySearchHistory
	result, err = testModifySearchHistory(user_Id)
	if err != nil || !result {
		log.Println("Failed to modify search history")
		return false, Data.UserConfig{}, err
	}
	return true, user_config, nil
}

// Userテストパート2
func TestApiUserPart2(userId string, site Data.WebSite, article Data.Article) (bool, error) {
	result, err := testModifyFavoriteSite(userId, site)
	if err != nil || !result {
		log.Println("Failed to modify favorite site")
		return false, err
	}
	result, err = testModifyFavoriteArticle(userId, article)
	if err != nil || !result {
		log.Println("Failed to modify favorite article")
		return false, err
	}
	result, err = testModifyAPIRequestLimit(userId,"Add")
	if err != nil || !result {
		log.Println("Failed to add api request limit")
		return false, err
	}
	result, err = testGetAPIRequestLimit(userId)
	if err != nil || !result {
		log.Println("Failed to get api request limit")
		return false, err
	}
	result, err = testModifyAPIRequestLimit(userId,"UnscopedDelete")
	if err != nil || !result {
		log.Println("Failed to unscoped delete api request limit")
		return false, err
	}
	return true, nil
}

func testModifyAPIRequestLimit(userId string,modify_type string) (bool, error) {
	request_type := "ModifyAPIRequestLimit"
	api_config := Data.ApiConfig{
		RefreshArticleInterval: 10,
	}
	api_config_json, _ := json.Marshal(api_config)
	api_config_str := string(api_config_json)
	result, err := SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType:          &request_type,
		UserId:               &userId,
		RequestArgumentJson1: &modify_type,
		RequestArgumentJson2: &api_config_str,
	})
	if err != nil || result != "Success ModifyAPIRequestLimit" {
		return false, err
	}
	log.Println("Success ModifyAPIRequestLimit")
	return true, nil
}

func genTestUserConfig() Data.UserConfig {
	return Data.UserConfig{
		UserName:     "test",
		UserUniqueID: "test-unique-id",
		ClientConfig: Data.ClientConfig{
			UiConfig: Data.UiConfig{
				ThemeColorValue:   677,
				DrawerMenuOpacity: 0.5,
				ArticleListFontSize: Data.UiResponsiveFontSize{
					Mobile: 10,
					Tablet: 12,
				},
				ThemeMode: "light",
				ArticleDetailFontSize: Data.UiResponsiveFontSize{
					Mobile: 10,
					Tablet: 12,
				},
			},
		},
		AccountType: "free",
		Country:     "Japan",
	}
}

func testRegisterUser(userId string) (bool, error) {
	request_type := "RegisterUser"
	user_config := genTestUserConfig()
	user_config_json, _ := json.Marshal(user_config)
	user_config_json_str := string(user_config_json)
	result, err := SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType:          &request_type,
		UserId:               &userId,
		RequestArgumentJson1: &user_config_json_str,
	})
	if err != nil || result != "Success RegisterUser" {
		log.Println("Failed to register user")
		return false, err
	}
	log.Println("Success to register user")
	return true, nil
}

// ConfigSync
func testConfigSync(userId string) (bool, error) {
	request_type := "ConfigSync"
	result, err := SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType: &request_type,
		UserId:      &userId,
	})
	if err != nil {
		log.Println("Failed to config sync")
		return false, err
	}
	// ApiConfigをパースする
	var config_sync_response = Data.ConfigSyncResponse{}
	err = json.Unmarshal([]byte(result), &config_sync_response)
	if err != nil || config_sync_response.UserConfig.ClientConfig.UiConfig.ThemeColorValue != 677 {
		log.Println("Failed to config sync")
		return false, err
	}
	log.Println("Success to config sync")
	return true, nil
}

// ReportReadActivity
func testReportReadActivity(userId string) (bool, error) {
	request_type := "ReportReadActivity"
	read_history := Data.ReadHistory{
		Link:           "https://www.google.com",
		AccessAt:       time.Now().Format(time.RFC3339),
		AccessPlatform: "web",
	}
	read_history_json, _ := json.Marshal(read_history)
	read_history_json_str := string(read_history_json)
	result, err := SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType:          &request_type,
		UserId:               &userId,
		RequestArgumentJson1: &read_history_json_str,
	})
	if err != nil || result != "Success ReportReadActivity" {
		log.Println("Failed to report read activity")
		return false, err
	}
	log.Println("Success to report read activity")
	return true, nil
}

// UpdateConfig
func testUpdateConfig(userId string) (bool, Data.UserConfig, error) {
	request_type := "UpdateConfig"
	user_config := genTestUserConfig()
	user_config.AccountType = "premium"
	user_config_json, _ := json.Marshal(user_config)
	user_config_json_str := string(user_config_json)
	result, err := SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType:          &request_type,
		UserId:               &userId,
		RequestArgumentJson1: &user_config_json_str,
	})
	if err != nil || result != "Success UpdateConfig" {
		log.Println("Failed to update config")
		return false, Data.UserConfig{}, err
	}
	// ユーザー設定を取得して、更新されているか確認する
	request_type = "ConfigSync"
	result, err = SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType: &request_type,
		UserId:      &userId,
	})
	if err != nil {
		log.Println("Failed to (update config):config sync")
		return false, Data.UserConfig{}, err
	}
	// UserConfigをパースする
	var config_sync_response = Data.ConfigSyncResponse{}
	err = json.Unmarshal([]byte(result), &config_sync_response)
	if err != nil || config_sync_response.UserConfig.AccountType != "premium" {
		log.Println("Failed to update config")
		return false, Data.UserConfig{}, err
	}
	log.Println("Success to update config")
	return false, config_sync_response.UserConfig, err
}

// ModifySearchHistory
func testModifySearchHistory(userId string) (bool, error) {
	req_type := "ModifySearchHistory"
	search_history := Data.SearchHistory{
		SearchWord: "test",
		SearchAt:   time.Now().Format(time.RFC3339),
	}
	search_history_json, _ := json.Marshal(search_history)
	search_history_json_str := string(search_history_json)
	result, err := SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType:          &req_type,
		UserId:               &userId,
		RequestArgumentJson1: &search_history_json_str,
	})
	expected_search_history_json, _ := json.Marshal([]string{"test"})
	if err != nil || result != string(expected_search_history_json) {
		log.Println("Failed to modify search history")
		return false, err
	}
	log.Println("Success to modify search history")
	return true, nil
}

// ModifyFavoriteSite
func testModifyFavoriteSite(userId string, favorite_site Data.WebSite) (bool, error) {
	request_type := "ModifyFavoriteSite"
	favorite_site_json, _ := json.Marshal(favorite_site)
	favorite_site_json_str := string(favorite_site_json)
	is_subscribe_json, _ := json.Marshal(true)
	is_subscribe_json_str := string(is_subscribe_json)
	result, err := SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType:          &request_type,
		UserId:               &userId,
		RequestArgumentJson1: &favorite_site_json_str,
		RequestArgumentJson2: &is_subscribe_json_str,
	})
	if err != nil || result != "Success ModifyFavoriteSite" {
		log.Println("Failed to modify favorite site")
		return false, err
	}
	log.Println("Success to modify favorite site")
	return true, nil
}

// ModifyFavoriteArticle
func testModifyFavoriteArticle(userId string, article Data.Article) (bool, error) {
	request_type := "ModifyFavoriteArticle"
	article_json, err := json.Marshal(article)
	if err != nil {
		log.Println("Failed to marshal article")
		return false, err
	}
	article_json_str := string(article_json)
	is_add_or_remove_json, err := json.Marshal(true)
	if err != nil {
		log.Println("Failed to marshal is add or remove")
		return false, err
	}
	is_add_or_remove_json_str := string(is_add_or_remove_json)
	result, err := SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType:          &request_type,
		UserId:               &userId,
		RequestArgumentJson1: &article_json_str,
		RequestArgumentJson2: &is_add_or_remove_json_str,
	})
	if err != nil || result != "Success ModifyFavoriteArticle" {
		log.Println("Failed to modify favorite article")
		return false, err
	}
	log.Println("Success to modify favorite article")
	return true, nil
}

// GetAPIRequestLimit
func testGetAPIRequestLimit(userId string) (bool, error) {
	request_type := "GetAPIRequestLimit"
	result, err := SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType: &request_type,
		UserId:      &userId,
	})
	if err != nil {
		log.Println("Failed to get api request limit")
		return false, err
	}
	expected_api_config := Data.ApiConfig{
		RefreshArticleInterval: 10,
	}
	// api_configをパースする
	var api_config = Data.ApiConfig{}
	err = json.Unmarshal([]byte(result), &api_config)
	if err != nil || api_config.RefreshArticleInterval != expected_api_config.RefreshArticleInterval {
		log.Println("Failed to get api request limit")
		return false, err
	}
	log.Println("Success to get api request limit")
	return true, nil
}

func SendUserRequest(request api_gen_code.PostUserJSONRequestBody) (string, error) {
	//  リクエストをjsonに変換
	request_post_json, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	// リクエストを作ったら、APIエンドポイントにリクエストを送る
	response, err := SendApiRequest(string(request_post_json), "user")
	if err != nil {
		return "", err
	}
	// resをData.APIResponseに変換
	var res = Data.APIResponse{}
	err = json.Unmarshal([]byte(*response.ResponseValue), &res)
	if err != nil {
		return "", err
	}
	return res.Value, nil
}
