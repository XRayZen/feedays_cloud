package APIFunction

import (
	"encoding/json"
	"read/Data"
)

func (s APIFunctions) SubscribeSite(access_ip string, user_id string, request_argument_json1 string, request_argument_json2 string) (string, error) {
	// サイトが登録されていなかったら登録処理をしてから購読を登録する
	// Feedの時間はRFC3339形式で返す
	var webSite Data.WebSite
	if err := json.Unmarshal([]byte(request_argument_json1), &webSite); err != nil {
		return "", err
	}
	var isSubscribe bool
	if err := json.Unmarshal([]byte(request_argument_json2), &isSubscribe); err != nil {
		return "", err
	}
	// サイトが購読されているか確認する
	if s.DBRepo.IsExistSite(webSite.SiteURL) {
		// 購読を登録する
		if err := s.DBRepo.SubscribeSite(user_id, webSite.SiteURL, isSubscribe); err != nil {
			return "", err
		}
	} else {
		// サイトのRSSを取得する
		articles, err := fetchRSSArticles(webSite.SiteRssURL)
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
	}
	return "", nil
}
