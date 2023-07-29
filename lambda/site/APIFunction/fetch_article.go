package APIFunction

import (
	"encoding/json"

	// "log"
	"site/Data"
	// "site/Repo"

	// "sort"
	"time"
)

func (s APIFunctions) FetchArticle(access_ip string, user_id string, request_argument_json1 string) (string, error) {
	// jsonからリクエストを変換する
	var request Data.FetchArticlesRequest
	err := json.Unmarshal([]byte(request_argument_json1), &request)
	if err != nil {
		return "", err
	}
	var response Data.FetchCloudFeedResponse
	// リクエストが新規なら最新記事を指定された件数で返す
	// リクエストが古いならクライアント側最古日時より古い記事を返す
	// リクエストが更新ならクライアント側更新日時より新しい記事を返す
	switch request.RequestType {
	case "Latest":
		articles, err := s.DBRepo.SearchSiteLatestArticle(request.SiteUrl, request.ReadCount)
		if err != nil {
			return "", err
		}
		response = Data.FetchCloudFeedResponse{
			Feeds:        articles,
			ResponseType: "success",
			Error:        "",
		}
	case "old":
		oldestModified, err := time.Parse(time.RFC3339, request.OldestModified)
		if err != nil {
			return "", err
		}
		// 指定された時間より古い記事を100件取得する
		articles, err := s.DBRepo.SearchArticlesByTimeAndOrder(request.SiteUrl, oldestModified, 100, false)
		if err != nil {
			return "", err
		}
		response = Data.FetchCloudFeedResponse{
			Feeds:        articles,
			ResponseType: "success",
			Error:        "",
		}
	case "Update":
		var newLastModified time.Time
		if request.LastModified == "" {
			newLastModified = time.Now()
		} else {
			newLastModified, err = time.Parse(time.RFC3339, request.LastModified)
			if err != nil {
				return "", err
			}
		}
		// サイトの記事更新日時を取得する
		articles, err := s.DBRepo.SearchArticlesByTimeAndOrder(request.SiteUrl, newLastModified, 100, true)
		if err != nil {
			return "", err
		}
		response = Data.FetchCloudFeedResponse{
			Feeds:        articles,
			ResponseType: "success",
			Error:        "",
		}
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(responseJson), nil
}
