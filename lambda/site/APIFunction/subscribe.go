package APIFunction

import (
	"encoding/json"
	"site/Data"
	// "site/Repo"
)

func (s APIFunctions) SubscribeSite(access_ip string, user_id string, request_argument_json1 string, request_argument_json2 string) (string, error) {
	// Feedの時間はRFC3339形式で返す
	var webSite Data.WebSite
	if err := json.Unmarshal([]byte(request_argument_json1), &webSite); err != nil {
		return "", err
	}
	var isSubscribe bool
	if err := json.Unmarshal([]byte(request_argument_json2), &isSubscribe); err != nil {
		return "", err
	}
	// サイトが登録されているか確認する
	if s.DBRepo.IsExistSite(webSite.SiteURL) {
		// 購読を登録・登録解除する
		if err := s.DBRepo.SubscribeSite(user_id, webSite.SiteURL, isSubscribe); err != nil {
			return "", err
		}
		return "Success Subscribe Site", nil
	} else {
		// サイトが登録されていなかったら登録処理をしてから購読を登録する
		// サイトのRSSを取得する
		no_image_articles, err := fetchRSSArticles(webSite.SiteRssURL)
		if err != nil {
			return "", err
		}
		// サイトのイメージURLを取得する
		articles, err := getArticleImageURLs(no_image_articles)
		if err != nil {
			return "", err
		}
		// サイトを登録する
		if err := s.DBRepo.RegisterSite(webSite, articles); err != nil {
			return "", err
		}
		// 購読を登録する
		if err := s.DBRepo.SubscribeSite(user_id, webSite.SiteURL, isSubscribe); err != nil {
			return "", err
		}
		return "Success Register Site", nil
	}
}
