package test_lambda_function

import (
	"encoding/json"
	"log"
	// "os/user"
	"test/Data"
	"test/api_gen_code"
	"time"
)

func TestApiUser() (bool, Data.UserConfig, error) {
	reqType := "GenUserID"
	requestResponse, err := SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType: &reqType,
	})
	if err != nil || len(requestResponse) < 10 {
		log.Println("Failed to get user id")
		return false, Data.UserConfig{}, err
	}
	log.Println("Success to get user id")
	userId := requestResponse
	// ユーザー登録
	result, err := testRegisterUser(userId)
	if err != nil || !result {
		log.Println("Failed to register user")
		return false, Data.UserConfig{}, err
	}
	// ConfigSync
	result, err = testConfigSync(userId)
	if err != nil || !result {
		log.Println("Failed to config sync")
		return false, Data.UserConfig{}, err
	}
	// ReportReadActivity
	result, err = testReportReadActivity(userId)
	if err != nil || !result {
		log.Println("Failed to report read activity")
		return false, Data.UserConfig{}, err
	}
	// UpdateConfig
	result, err = testUpdateConfig(userId)
	if err != nil || !result {
		log.Println("Failed to update config")
		return false, Data.UserConfig{}, err
	}
	// ModifySearchHistory
	result, err = testModifySearchHistory(userId)
	if err != nil || !result {
		log.Println("Failed to modify search history")
		return false, Data.UserConfig{}, err
	}

}

func genTestUserConfig() Data.UserConfig {
	return Data.UserConfig{
		UserName:     "test",
		UserUniqueID: "test-unique-id",
		ClientConfig: Data.ClientConfig{
			ApiConfig: Data.ApiConfig{
				RefreshArticleInterval: 10,
			},
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
	reqType := "RegisterUser"
	cfg := genTestUserConfig()
	cfgJson, err := json.Marshal(cfg)
	if err != nil {
		log.Println("Failed to marshal user config")
		return false, err
	}
	cfgStr := string(cfgJson)
	result, err := SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType:          &reqType,
		UserId:               &userId,
		RequestArgumentJson1: &cfgStr,
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
	reqType := "ConfigSync"
	result, err := SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType: &reqType,
		UserId:      &userId,
	})
	if err != nil {
		log.Println("Failed to config sync")
		return false, err
	}
	// ApiConfigをパースする
	var ConfigSyncResponse = Data.ConfigSyncResponse{}
	err = json.Unmarshal([]byte(result), &ConfigSyncResponse)
	if err != nil || ConfigSyncResponse.UserConfig.ClientConfig.ApiConfig.RefreshArticleInterval != 10 {
		log.Println("Failed to config sync")
		return false, err
	}
	log.Println("Success to config sync")
	return true, nil
}

// ReportReadActivity
func testReportReadActivity(userId string) (bool, error) {
	reqType := "ReportReadActivity"
	readHist := Data.ReadHistory{
		Link:           "https://www.google.com",
		AccessAt:       time.Now().Format(time.RFC3339),
		AccessPlatform: "web",
	}
	readHistJson, err := json.Marshal(readHist)
	if err != nil {
		log.Println("Failed to marshal read history")
		return false, err
	}
	readHistJsonStr := string(readHistJson)
	result, err := SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType:          &reqType,
		UserId:               &userId,
		RequestArgumentJson1: &readHistJsonStr,
	})
	if err != nil || result != "Success ReportReadActivity" {
		log.Println("Failed to report read activity")
		return false, err
	}
	log.Println("Success to report read activity")
	return true, nil
}

// UpdateConfig
func testUpdateConfig(userId string) (bool, error) {
	reqType := "UpdateConfig"
	cfg := genTestUserConfig()
	cfg.AccountType = "premium"
	cfgJson, err := json.Marshal(cfg)
	if err != nil {
		log.Println("Failed to marshal user config")
		return false, err
	}
	cfgStr := string(cfgJson)
	result, err := SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType:          &reqType,
		UserId:               &userId,
		RequestArgumentJson1: &cfgStr,
	})
	if err != nil || result != "Success UpdateConfig" {
		log.Println("Failed to update config")
		return false, err
	}
	// ユーザー設定を取得して、更新されているか確認する
	reqType = "ConfigSync"
	result, err = SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType: &reqType,
		UserId:      &userId,
	})
	// UserConfigをパースする
	var ConfigSyncResponse = Data.ConfigSyncResponse{}
	err = json.Unmarshal([]byte(result), &ConfigSyncResponse)
	if err != nil || ConfigSyncResponse.UserConfig.AccountType != "premium" {
		log.Println("Failed to update config")
		return false, err
	}
	log.Println("Success to update config")
	return true, nil
}

// ModifySearchHistory
func testModifySearchHistory(userId string) (bool, error) {
	reqType := "ModifySearchHistory"
	searchHist := Data.SearchHistory{
		SearchWord: "test",
		SearchAt:   time.Now().Format(time.RFC3339),
	}
	searchHistJson, err := json.Marshal(searchHist)
	if err != nil {
		log.Println("Failed to marshal search history")
		return false, err
	}
	searchHistJsonStr := string(searchHistJson)
	result, err := SendUserRequest(api_gen_code.PostUserJSONRequestBody{
		RequestType:          &reqType,
		UserId:               &userId,
		RequestArgumentJson1: &searchHistJsonStr,
	})
	searchHistoryJson, _ := json.Marshal([]string{"test"})
	if err != nil || result != string(searchHistoryJson) {
		log.Println("Failed to modify search history")
		return false, err
	}
	log.Println("Success to modify search history")
	return true, nil
}

// ModifyFavoriteSite



func SendUserRequest(request api_gen_code.PostUserJSONRequestBody) (string, error) {
	//  リクエストをjsonに変換
	requestPostJson, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	// リクエストを作ったら、APIエンドポイントにリクエストを送る
	response, err := SendApiRequest(string(requestPostJson), "user")
	if err != nil {
		return "", err
	}
	// resをData.APIResponseに変換
	var res = Data.APIResponse{}
	err = json.Unmarshal([]byte(*response.ResponseValue), &res)
	if err != nil {
		return "", err
	}
	return res.Message, nil
}
