package test_lambda_function

import (
	"encoding/json"
	"log"
	"net/http"
	// "os"
	"strings"
	"test/api_gen_code"
)

func LambdaApiTest() (bool, error) {
	// テストを実行する為にサービスを初期化する
	result, err := ServiceInit()
	if err != nil || !result {
		log.Println("ServiceInit: Failed ")
		if err != nil {
			log.Println("ServiceInit: Failed Detail Error :", err)
		}
		return false, err
	}
	// Userのテスト
	// ユーザーID生成・ユーザー登録・ ConfigSync ReportReadActivity UpdateConfig ModifySearchHistory
	result, user_cfg, err := TestApiUserPart1()
	if err != nil || !result {
		log.Println("TestApiUserPart1: Failed")
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
	result,err = DeleteTestData(user_cfg.UserUniqueID, sites[0].SiteURL)
	if err != nil || !result {
		log.Println("DeleteTestData: Failed : ", err)
		return false, err
	}
	log.Println("ALL TESTS ARE SUCCESSFUL")
	return true, nil
}

// リクエストを送信
func SendApiRequest(request_post_json string, path string) (api_gen_code.APIResponse, error) {
	// APIエンドポイントにリクエストを送る
	end_point := "https://bkq8lpslz8.execute-api.us-east-1.amazonaws.com/develop"
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
		log.Println("Failed Decode: ", err)
		return api_gen_code.APIResponse{}, err
	}
	return response, nil
}
