package test_lambda_function

import (
	"encoding/json"
	"log"
	"test/Data"
	"test/api_gen_code"
	"time"
	// "github.com/aws/aws-lambda-go/events"
)

func SendSearchRequest(request Data.ApiSearchRequest) (Data.SearchResult, error) {
	// リクエストをjsonに変換する
	requestJson, err := json.Marshal(request)
	if err != nil {
		return Data.SearchResult{}, err
	}
	requestJsonStr := string(requestJson)
	requestTypeStr := "Search"
	userID := "0000"
	PostRequestSite := api_gen_code.PostSiteJSONRequestBody{
		RequestType:          &requestTypeStr,
		UserId:               &userID,
		RequestArgumentJson1: &requestJsonStr,
	}
	// リクエストをjsonに変換する
	requestPostJson, err := json.Marshal(PostRequestSite)
	if err != nil {
		return Data.SearchResult{}, err
	}
	// リクエストを作ったら、APIエンドポイントにリクエストを送る
	response, err := SendApiRequest(string(requestPostJson), "site")
	if err != nil {
		return Data.SearchResult{}, err
	}
	// レスポンスをパースする
	var result Data.SearchResult
	err = json.Unmarshal([]byte(*response.ResponseValue), &result)
	if err != nil {
		return Data.SearchResult{}, err
	}
	return result, nil
}

func SendSubscribeSiteRequest(requestWebSite Data.WebSite, isSubscribe bool) (string, error) {
	// サイトとisSubscribeをjsonに変換する
	requestJson, err := json.Marshal(requestWebSite)
	if err != nil {
		return "", err
	}
	requestJsonStr := string(requestJson)
	isSubscribeJson, err := json.Marshal(isSubscribe)
	if err != nil {
		return "", err
	}
	isSubscribeJsonStr := string(isSubscribeJson)
	reqType := "SubscribeSite"
	userId := "000"
	request := api_gen_code.PostSiteJSONRequestBody{
		RequestType:          &reqType,
		UserId:               &userId,
		RequestArgumentJson1: &requestJsonStr,
		RequestArgumentJson2: &isSubscribeJsonStr,
	}
	// リクエストをjsonに変換する
	requestPostJson, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	// リクエストを作ったら、APIエンドポイントにリクエストを送る
	response, err := SendApiRequest(string(requestPostJson), "site")
	if err != nil {
		return "", err
	}
	return *response.ResponseValue, nil
}

func SendFetchArticleRequest(request Data.FetchArticlesRequest) (string, error) {
	// リクエストをjsonに変換する
	requestJson, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	requestJsonStr := string(requestJson)
	reqType := "FetchArticle"
	userId := "000"
	requestBody := api_gen_code.PostSiteJSONRequestBody{
		RequestType:          &reqType,
		UserId:               &userId,
		RequestArgumentJson1: &requestJsonStr,
	}
	// リクエストをjsonに変換する
	requestPostJson, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	// リクエストを作ったら、APIエンドポイントにリクエストを送る
	response, err := SendApiRequest(string(requestPostJson), "site")
	if err != nil {
		return "", err
	}
	return *response.ResponseValue, nil
}

// 新規サイトを登録して登録できたかどうかを確かめる
func TestApiSiteSearch() (bool, []Data.WebSite, error) {
	// ギガジンをURL検索してちゃんと登録に成功したかどうかを確かめる
	request := Data.ApiSearchRequest{
		SearchType: "URL",
		Word:       "https://gigazine.net/",
	}
	result, err := SendSearchRequest(request)
	if err != nil {
		return false, nil, err
	}
	// サイト名がギガジンであることを確かめる
	if result.Websites[0].SiteName != "GIGAZINE" {
		return false, nil, err
	}
	// テストが成功したことをログに出力する
	log.Println("TestApiSiteNewSiteSearch: Success")
	request = Data.ApiSearchRequest{
		SearchType: "Keyword",
		Word:       "AI",
	}
	result, err = SendSearchRequest(request)
	if err != nil {
		return false, nil, err
	}
	// ちゃんと記事が返ってきているかどうかを確かめる
	if len(result.Articles) == 0 {
		return false, nil, err
	}
	// テストが成功したことをログに出力する
	log.Println("TestApiSiteKeywordSearch: Success")
	request = Data.ApiSearchRequest{
		SearchType: "SiteName",
		Word:       "GIGAZINE",
	}
	result, err = SendSearchRequest(request)
	if err != nil {
		return false, nil, err
	}
	// サイト名がギガジンであることを確かめる
	if result.Websites[0].SiteName != "GIGAZINE" {
		return false, nil, err
	}
	// テストが成功したことをログに出力する
	log.Println("TestApiSiteSiteNameSearch: Success")
	return true, result.Websites, nil
}

// 返ってきたサイトを購読して購読できたかどうかを確かめる
func TestApiSiteSubscribe(site Data.WebSite) (bool, error) {
	res, err := SendSubscribeSiteRequest(site, true)
	if err != nil {
		return false, err
	}
	if res != "Success Subscribe Site" {
		return false, err
	}
	log.Println("TestApiSiteSubscribe: Success")
	res, err = SendSubscribeSiteRequest(site, false)
	if err != nil {
		return false, err
	}
	if res != "Success UnSubscribe Site" {
		return false, err
	}
	log.Println("TestApiSiteUnSubscribe: Success")
	return true, nil
}

// サイトの記事を取得する
func TestApiSiteFetchArticles(site Data.WebSite) (bool, error) {
	// サイトの最新記事を取得する
	requestFetchArticleByLatest := Data.FetchArticlesRequest{
		SiteUrl:     site.SiteURL,
		RequestType: "Latest",
		ReadCount:   10,
	}
	res, err := SendFetchArticleRequest(requestFetchArticleByLatest)
	if err != nil {
		return false, err
	}
	var result Data.FetchArticleResponse
	err = json.Unmarshal([]byte(res), &result)
	if err != nil {
		return false, err
	}
	// 評価する
	if len(result.Articles) == 0 {
		return false, err
	}
	// テストが成功したことをログに出力する
	log.Println("TestApiSiteFetchArticlesByLatest: Success")
	// 指定された時間より古い記事を100件取得する
	// 1日前の時間を取得する
	now := time.Now()
	oneDayAgo := now.AddDate(0, 0, -1)
	// 1日前の時間をRFC3339文字列に変換する
	oneDayAgoStr := oneDayAgo.Format(time.RFC3339)
	requestOlder := Data.FetchArticlesRequest{
		SiteUrl:     site.SiteURL,
		RequestType: "Older",
		OldestModified: oneDayAgoStr,
	}
	res, err = SendFetchArticleRequest(requestOlder)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal([]byte(res), &result)
	if err != nil {
		return false, err
	}
	// 評価する
	if len(result.Articles) == 0 {
		return false, err
	}
	// テストが成功したことをログに出力する
	log.Println("TestApiSiteFetchArticlesByOlder: Success")
	// 更新をテストだが、ここから更新はテスト出来ない 該当関数内で十分テストされている
	return true, nil
}
