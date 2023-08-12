package test_lambda_function

import (
	"encoding/json"
	"log"
	"test/Data"
	"test/api_gen_code"
	"time"
	// "github.com/aws/aws-lambda-go/events"
)

// 新規サイトを登録して登録できたかどうかを確かめる
func TestApiSitePart1(userId string) (bool, []Data.WebSite, error) {
	// URL検索してちゃんと登録に成功したかどうかを確かめる
	request := Data.ApiSearchRequest{
		SearchType: "URL",
		Word:       "https://gigazine.net/",
	}
	result, err := SendSearchRequest(request, userId)
	if err != nil || result.Websites[0].SiteName != "GIGAZINE" {
		log.Println("TestApiSiteUrlSearch: Failed")
		return false, nil, err
	}
	log.Println("TestApiSiteNewSiteSearch: Success")
	// キーワードで検索してちゃんと登録に成功したかどうかを確かめる
	request = Data.ApiSearchRequest{
		SearchType: "Keyword",
		Word:       "AI",
	}
	result, err = SendSearchRequest(request, userId)
	// ちゃんと記事が返ってきているかどうかを確かめる
	if err != nil || len(result.Articles) == 0 {
		log.Println("TestApiSiteKeywordSearch: Failed")
		return false, nil, err
	}
	log.Println("TestApiSiteKeywordSearch: Success")
	// サイト名で検索してちゃんと登録に成功したかどうかを確かめる
	request = Data.ApiSearchRequest{
		SearchType: "SiteName",
		Word:       "GIGAZINE",
	}
	result, err = SendSearchRequest(request, userId)
	if err != nil || result.Websites[0].SiteName != "GIGAZINE" {
		return false, nil, err
	}
	log.Println("TestApiSiteSiteNameSearch: Success")
	return true, result.Websites, nil
}

// ExploreCategoriesを追加してサイトにカテゴリー追加してカテゴリー検索できるかどうかを確かめる
func TestApiSitePart2(site Data.WebSite, userID string) (bool, error) {
	req_type := "ModifyExploreCategory"
	explore_category := Data.ExploreCategory{
		CategoryName:        "IT",
		CategoryDescription: "IT",
		CategoryCountry:     "Japan",
	}
	explore_category_json, _ := json.Marshal(explore_category)
	explore_category_json_str := string(explore_category_json)
	is_add_or_remove_json, _ := json.Marshal(true)
	is_add_or_remove_json_str := string(is_add_or_remove_json)
	request := api_gen_code.PostSiteJSONRequestBody{
		RequestType:          &req_type,
		UserId:               &userID,
		RequestArgumentJson1: &explore_category_json_str,
		RequestArgumentJson2: &is_add_or_remove_json_str,
	}
	// リクエストをjsonに変換する
	request_post_json, _ := json.Marshal(request)
	// リクエストを作ったら、APIエンドポイントにリクエストを送る
	response, err := SendApiRequest(string(request_post_json), "site")
	if err != nil {
		return false, err
	}
	if *response.ResponseValue != "Success ModifyExploreCategory" {
		return false, err
	}
	log.Println("TestApiSiteModifyExploreCategory: Success")
	// サイトにカテゴリーを追加する
	req_type = "ChangeSiteCategory"
	site.SiteCategory = "IT"
	site_url_json, _ := json.Marshal(site.SiteURL)
	site_url_json_str := string(site_url_json)
	category_name := "IT"
	category_name_json, _ := json.Marshal(category_name)
	category_name_json_str := string(category_name_json)
	request = api_gen_code.PostSiteJSONRequestBody{
		RequestType:          &req_type,
		UserId:               &userID,
		RequestArgumentJson1: &site_url_json_str,
		RequestArgumentJson2: &category_name_json_str,
	}
	// リクエストをjsonに変換する
	request_post_json, _ = json.Marshal(request)
	// リクエストを作ったら、APIエンドポイントにリクエストを送る
	response, err = SendApiRequest(string(request_post_json), "site")
	if err != nil {
		return false, err
	}
	if *response.ResponseValue != "Success ChangeSiteCategory" {
		return false, err
	}
	log.Println("TestApiSiteChangeSiteCategory: Success")
	// カテゴリー検索をする
	request_search_by_category := Data.ApiSearchRequest{
		SearchType: "Category",
		Word:       "IT",
		UserID:     userID,
	}
	result, err := SendSearchRequest(request_search_by_category, userID)
	if err != nil {
		return false, err
	}
	if len(result.Websites) == 0 {
		return false, err
	}
	log.Println("TestApiSiteSearchByCategory: Success")
	return true, nil
}

// 返ってきたサイトを購読して購読できたかどうかを確かめる
func TestApiSitePart3(site Data.WebSite, userID string) (bool, error) {
	res, err := SendSubscribeSiteRequest(site, true, userID)
	if err != nil {
		return false, err
	}
	if res != "Success Subscribe Site" {
		return false, err
	}
	log.Println("TestApiSiteSubscribe: Success")
	res, err = SendSubscribeSiteRequest(site, false, userID)
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
func TestApiSitePart4(site Data.WebSite, userID string) (bool, Data.Article, error) {
	// サイトの最新記事を取得する
	request_fetch_article_by_latest := Data.FetchArticlesRequest{
		SiteUrl:     site.SiteURL,
		RequestType: "Latest",
		ReadCount:   10,
	}
	res, err := SendFetchArticleRequest(request_fetch_article_by_latest, userID)
	if err != nil {
		return false, Data.Article{}, err
	}
	var result_fetch_article_by_latest Data.FetchArticleResponse
	err = json.Unmarshal([]byte(res), &result_fetch_article_by_latest)
	if err != nil || len(result_fetch_article_by_latest.Articles) == 0 {
		log.Println("TestApiSiteFetchArticlesByLatest: Failed")
		return false, Data.Article{}, err
	}
	// テストが成功したことをログに出力する
	log.Println("TestApiSiteFetchArticlesByLatest: Success")
	// 指定された時間より古い記事を100件取得する
	// 1日前の時間を取得する
	now := time.Now()
	one_day_ago := now.AddDate(0, 0, -1)
	// 1日前の時間をRFC3339文字列に変換する
	one_day_ago_str := one_day_ago.Format(time.RFC3339)
	request_older := Data.FetchArticlesRequest{
		SiteUrl:        site.SiteURL,
		RequestType:    "Older",
		OldestModified: one_day_ago_str,
	}
	res, err = SendFetchArticleRequest(request_older, userID)
	if err != nil {
		return false, Data.Article{}, err
	}
	var result Data.FetchArticleResponse
	err = json.Unmarshal([]byte(res), &result)
	if err != nil || len(result.Articles) == 0 {
		return false, Data.Article{}, err
	}
	// テストが成功したことをログに出力する
	log.Println("TestApiSiteFetchArticlesByOlder: Success")
	// 更新をテストだが、ここから更新はテスト出来ない 該当関数内で十分テストされている
	return true, result_fetch_article_by_latest.Articles[0], nil
}

func SendSearchRequest(request Data.ApiSearchRequest, userID string) (Data.SearchResult, error) {
	// リクエストをjsonに変換する
	request_json, err := json.Marshal(request)
	if err != nil {
		return Data.SearchResult{}, err
	}
	request_json_str := string(request_json)
	request_type_str := "Search"
	post_request_site := api_gen_code.PostSiteJSONRequestBody{
		RequestType:          &request_type_str,
		UserId:               &userID,
		RequestArgumentJson1: &request_json_str,
	}
	// リクエストをjsonに変換する
	request_post_json, err := json.Marshal(post_request_site)
	if err != nil {
		return Data.SearchResult{}, err
	}
	// リクエストを作ったら、APIエンドポイントにリクエストを送る
	response, err := SendApiRequest(string(request_post_json), "site")
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

func SendSubscribeSiteRequest(requestWebSite Data.WebSite, isSubscribe bool, userID string) (string, error) {
	// サイトとisSubscribeをjsonに変換する
	request_json, err := json.Marshal(requestWebSite)
	if err != nil {
		return "", err
	}
	request_json_str := string(request_json)
	is_subscribe_json, err := json.Marshal(isSubscribe)
	if err != nil {
		return "", err
	}
	is_subscribe_json_str := string(is_subscribe_json)
	request_type := "SubscribeSite"
	post_request := api_gen_code.PostSiteJSONRequestBody{
		RequestType:          &request_type,
		UserId:               &userID,
		RequestArgumentJson1: &request_json_str,
		RequestArgumentJson2: &is_subscribe_json_str,
	}
	// リクエストをjsonに変換する
	request_post_json, err := json.Marshal(post_request)
	if err != nil {
		return "", err
	}
	// リクエストを作ったら、APIエンドポイントにリクエストを送る
	response, err := SendApiRequest(string(request_post_json), "site")
	if err != nil {
		return "", err
	}
	return *response.ResponseValue, nil
}

func SendFetchArticleRequest(request Data.FetchArticlesRequest, userID string) (string, error) {
	// リクエストをjsonに変換する
	request_json, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	request_json_str := string(request_json)
	request_type := "FetchArticle"
	request_body := api_gen_code.PostSiteJSONRequestBody{
		RequestType:          &request_type,
		UserId:               &userID,
		RequestArgumentJson1: &request_json_str,
	}
	// リクエストをjsonに変換する
	request_post_json, err := json.Marshal(request_body)
	if err != nil {
		return "", err
	}
	// リクエストを作ったら、APIエンドポイントにリクエストを送る
	response, err := SendApiRequest(string(request_post_json), "site")
	if err != nil {
		return "", err
	}
	return *response.ResponseValue, nil
}
