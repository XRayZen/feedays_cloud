package test_lambda_function

import (
	"encoding/json"
	"log"
	"test/Data"

	"github.com/aws/aws-lambda-go/events"
)

// リクエストを送信
func SendRequest(reqStr string,apiRequestType string)(events.APIGatewayProxyResponse, error) {
	// OpenAPIのコードでリクエストを作ってAPIエンドポイントにリクエストを送る
}

// 新規サイトを登録して登録できたかどうかを確かめる
func TestApiSiteNewSiteSearch() (bool, error) {
	// ギガジンをURL検索してちゃんと登録に成功したかどうかを確かめる
	request := Data.ApiSearchRequest{
		SearchType: "URL",
		Word:       "https://gigazine.net/",
	}
	// リクエストをjsonに変換する
	requestJson, err := json.Marshal(request)
	if err != nil {
		return false, err
	}
	// リクエストを作ったら、APIエンドポイントにリクエストを送る
	response, err := SendRequest(string(requestJson))
	if err != nil {
		return false, err
	}
	// レスポンスをjsonに変換する
	var result = Data.SearchResult{}
	if err := json.Unmarshal([]byte(response.Body), &result); err != nil {
		return false, err
	}
	// サイト名がギガジンであることを確かめる
	if result.Websites[0].SiteName != "GIGAZINE" {
		return false, nil
	}
	// テストが成功したことをログに出力する
	log.Println("TestApiSiteNewSiteSearch: Success")
	return true, nil
}

// 記事をキーワード検索してちゃんと成功したかどうかを確かめる
func TestApiSiteKeywordSearch() (bool, error) {
	request := Data.ApiSearchRequest{
		SearchType: "Keyword",
		Word:       "AI",
	}
	// リクエストをjsonに変換する
	requestJson, err := json.Marshal(request)
	if err != nil {
		return false, err
	}
	// リクエストを作ったら、APIエンドポイントにリクエストを送る
	response, err := SendRequest(string(requestJson))
	if err != nil {
		return false, err
	}
	// レスポンスをjsonに変換する
	var result = Data.SearchResult{}
	if err := json.Unmarshal([]byte(response.Body), &result); err != nil {
		return false, err
	}
	// ちゃんと記事が返ってきているかどうかを確かめる
	if len(result.Articles) == 0 {
		return false, nil
	}
	// テストが成功したことをログに出力する
	log.Println("TestApiSiteKeywordSearch: Success")
	return true, nil
}

// ギガジンをサイト名検索してちゃんと成功したかどうかを確かめる
func TestApiSiteSiteNameSearch() (bool, error) {
	request := Data.ApiSearchRequest{
		SearchType: "SiteName",
		Word:       "GIGAZINE",
	}
	// リクエストをjsonに変換する
	requestJson, err := json.Marshal(request)
	if err != nil {
		return false, err
	}
	// リクエストを作ったら、APIエンドポイントにリクエストを送る
	response, err := SendRequest(string(requestJson))
	if err != nil {
		return false, err
	}
	// レスポンスをjsonに変換する
	var result = Data.SearchResult{}
	if err := json.Unmarshal([]byte(response.Body), &result); err != nil {
		return false, err
	}
	// サイト名がギガジンであることを確かめる
	if result.Websites[0].SiteName != "GIGAZINE" {
		return false, nil
	}
	// テストが成功したことをログに出力する
	log.Println("TestApiSiteSiteNameSearch: Success")
	return true, nil
}
