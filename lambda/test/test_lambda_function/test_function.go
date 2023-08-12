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

func LambdaApiTest() (bool, error) {
	// Userのテスト
	// ユーザーID生成・ユーザー登録・ ConfigSync ReportReadActivity UpdateConfig ModifySearchHistory
	result, user_cfg, err := TestApiUserPart1()
	if err != nil || !result {
		return false, err
	}
	// Siteのテスト
	// サイトを登録して検索するまで
	result, sites, err := TestApiSitePart1(user_cfg.UserUniqueID)
	if err != nil || !result {
		log.Println("TestApiSiteSearch: Failed")
		return false, err
	}
	// ExploreCategoriesのテスト
	// カテゴリーを追加する
	result, err = TestApiSitePart2(sites[0], user_cfg.UserUniqueID)
	if err != nil || !result {
		log.Println("TestApiSiteModifyExploreCategory: Failed")
		return false, err
	}
	result, article, err := TestApiSitePart4(sites[0], user_cfg.UserUniqueID)
	if err != nil || !result {
		log.Println("TestApiSiteFetchArticles: Failed")
		return false, err
	}
	// Userのテスト
	result, err = TestApiUserPart2(user_cfg.UserUniqueID, sites[0], article)
	if err != nil || !result {
		log.Println("TestApiUserPart2: Failed")
		return false, err
	}
	// 購読する
	result, err = TestApiSitePart3(sites[0], user_cfg.UserUniqueID)
	if err != nil || !result {
		log.Println("TestApiSiteSubscribe: Failed")
		return false, err
	}
	// 次はReadのテスト
	result,user_cfg, err = TestApiRead(user_cfg.UserUniqueID)
	if err != nil || !result {
		log.Println("TestApiRead: Failed")
		return false, err
	}
	// テストが全て成功したら削除する
	
	return true, nil
}

// リクエストを送信
func SendApiRequest(request_post_json string, path string) (api_gen_code.APIResponse, error) {
	// APIエンドポイントにリクエストを送る
	end_point := os.Getenv("API_ENDPOINT")
	// そこに/siteなどをつける
	end_point = end_point + "/" + path
	// HTTPポストリクエストを送る
	res, err := http.Post(end_point, "application/json", strings.NewReader(request_post_json))
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
