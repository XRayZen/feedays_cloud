package APIFunction

import (
	"encoding/json"
	"errors"

	// "log"
	"read/Data"
	"read/Repo"

	// "sort"
	"time"
)

func (s APIFunctions) FetchCloudFeed(access_ip string, user_id string, request_argument_json1 string) (string, error) {
	// jsonからURLを取得する
	var request Data.FetchFeedRequest
	err := json.Unmarshal([]byte(request_argument_json1), &request)
	if err != nil {
		return "", err
	}
	clientLastModified, err := time.Parse(time.RFC3339, request.LastModified)
	if err != nil {
		return "", err
	}
	articles, err := refreshSiteArticles(s.DBRepo, request.SiteUrl, request.IntervalMinutes, clientLastModified)
	if err != nil {
		return "", err
	}
	var response Data.FetchCloudFeedResponse
	response.Feeds = articles
	response.ResponseType = "success"
	responseJson, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(responseJson), nil
}

// 指定したサイトURLの記事が指定された間隔より古くなっていたら更新する
// まだ鮮度があるのなら更新せずそのままクライアント側更新日時より新しい記事を返す
func refreshSiteArticles(repo Repo.DBRepository, siteUrl string, intervalMinutes int, clientLastModified time.Time) ([]Data.Article, error) {
	if repo.IsExistSite(siteUrl) {
		// サイトの記事更新日時を取得する
		lastModified, err := repo.FetchSiteLastModified(siteUrl)
		if err != nil {
			return nil, err
		}
		// 記事更新日時にIntervalMinutesを足した更新期限日時に現時間が過ぎていたら記事更新をする
		if isUpdateExpired(lastModified, intervalMinutes) {
			// サイトを取得する
			site, err := repo.FetchSite(siteUrl)
			if err != nil {
				return nil, err
			}
			// サイトの記事を取得する
			articles, err := fetchRSSArticles(site.SiteRssURL)
			if err != nil {
				return nil, err
			}
			// サイトの記事をDBに登録する
			err = repo.UpdateArticles(siteUrl, articles)
			if err != nil {
				return nil, err
			}
		}
		// クライアント側更新日時より新しい記事を返す
		articles, err := repo.GetArticlesByTme(siteUrl, clientLastModified)
		if err != nil {
			return nil, err
		}
		return articles, nil
	}
	return nil, errors.New("Not Found WebSite")
}

// 記事更新日時にIntervalMinutesを足した更新期限日時を現時間が過ぎていたらtrueを返す
func isUpdateExpired(lastModified time.Time, intervalMinutes int) bool {
	// 現時間を取得する
	now_time := time.Now()
	// 記事更新日時にIntervalMinutesを足した更新期限日時を取得する
	update_expired_time := lastModified.Add(time.Minute * time.Duration(intervalMinutes))
	// 更新期限日時が現時間より過ぎていたらtrueを返す
	return update_expired_time.Before(now_time)
}
