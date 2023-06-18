package Repo

import (
	// "log"
	"read/Data"
	"time"
)

type MockDBRepo struct {
	articles []Data.Article
}

// readで使う
func (s MockDBRepo) GetUserInfo(userID string) (resUserInfo Data.UserInfo, err error) {
	return Data.UserInfo{}, nil
}

func (s MockDBRepo) GetExploreCategories(userID string, country string) (resExp Data.ExploreCategories, err error) {
	return Data.ExploreCategories{
		CategoryName: "CategoryName",
	}, nil
}

func (s MockDBRepo) GetRanking(userID string, country string) (resRanking Data.Ranking, err error) {
	return Data.Ranking{
		RankingName: "RankingName",
	}, nil
}

// heavyで使う
func (s MockDBRepo) IsExistSite(site_url string) bool {
	switch site_url {
	case "https://automaton-media.com/":
		return true
	case "https://gigazine.net/":
		return true
	}
	return false
}

func (s MockDBRepo) IsSubscribeSite(user_id string, site_url string) bool {
	return true
}

func (s MockDBRepo) FetchSite(site_url string) (Data.WebSite, error) {
	// テスト用にダミーのWebSiteを返す
	switch site_url {
	case "https://automaton-media.com/":
		// サイトが存在していて更新期限を過ぎている場合
		// 30分前の日時を返す
		lastModifiedTime := time.Now().Add(-time.Minute * 30)
		return Data.WebSite{
			SiteURL:      "https://automaton-media.com/",
			SiteRssURL:   "https://automaton-media.com/feed/",
			SiteName:     "AUTOMATON",
			LastModified: lastModifiedTime.Format(time.RFC3339),
		}, nil
	case "https://gigazine.net/":
		// サイトが存在していて更新期限を過ぎていない場合
		// 現在の日時
		lastModifiedTime := time.Now()
		return Data.WebSite{
			SiteURL:      "https://gigazine.net/",
			SiteRssURL:   "https://gigazine.net/news/rss_2.0/",
			SiteName:     "GIGAZINE",
			LastModified: lastModifiedTime.Format(time.RFC3339),
		}, nil
	}
	return Data.WebSite{}, nil
}

func (s MockDBRepo) FetchSiteLastModified(site_url string) (time.Time, error) {
	// テスト用にダミーのWebSiteを返す
	switch site_url {
	case "https://automaton-media.com/":
		// サイトが存在していて更新期限を過ぎている場合
		// 二時間前の日時を返す
		lastModifiedTime := time.Now().Add(-time.Hour * 2)
		return lastModifiedTime, nil
	case "https://gigazine.net/":
		// サイトが存在していて更新期限を過ぎていない場合
		// 現在の日時
		lastModifiedTime := time.Now()
		return lastModifiedTime, nil
	}
	return time.Now(), nil
}

func (s MockDBRepo) RegisterSite(site Data.WebSite, articles []Data.Article) error {
	return nil
}

func (s MockDBRepo) SearchArticlesByKeyword(keyword string) ([]Data.Article, error) {
	if keyword == "Found" {
		return []Data.Article{
			{
				Title:        "Found",
				Link:         "https://example.com",
				Site:         "https://example.com",
				LastModified: "2021-01-01T00:00:00+09:00",
			},
		}, nil
	}
	return nil, nil
}

func (s MockDBRepo) SearchArticlesByTime(siteUrl string, lastModified time.Time) ([]Data.Article, error) {
	// 更新日時より新しい記事を返す
	var articles []Data.Article
	for _, article := range mockArticles {
		articleTime, _ := time.Parse(time.RFC3339, article.LastModified)
		// articleTimeを数値に変換
		articleTimeUnix := articleTime.Unix()
		// lastModifiedを数値に変換
		lastModifiedUnix := lastModified.Unix()
		// log.Println("articleTimeUnix > lastModifiedUnix: ", articleTimeUnix > lastModifiedUnix)
		// lastModifiedよりarticleTimeが新しい場合は追加する
		if articleTimeUnix > lastModifiedUnix {
			articles = append(articles, article)
		}
	}
	return articles, nil
}

func (s MockDBRepo) SearchSiteByName(siteName string) ([]Data.WebSite, error) {
	switch siteName {
	case "Found":
		return []Data.WebSite{
			{
				SiteURL:      "https://example.com",
				SiteRssURL:   "https://example.com",
				SiteName:     "Found",
				LastModified: "2021-01-01T00:00:00+09:00",
			},
		}, nil
	}
	return nil, nil
}

func (s MockDBRepo) SearchArticlesBySite(siteUrl string) ([]Data.Article, error) {
	return nil, nil
}

// モック用変数
var mockArticles = []Data.Article{}

func (s MockDBRepo) UpdateArticles(siteUrl string, articles []Data.Article) error {
	mockArticles = append(mockArticles, articles...)
	return nil
}

func (s MockDBRepo) SubscribeSite(user_id string, siteUrl string, is_subscribe bool) error {
	return nil
}

// バッチ処理用
func (s MockDBRepo) FetchAllSites() ([]Data.WebSite, error) {
	return []Data.WebSite{}, nil
}

func (s MockDBRepo) FetchAllHistories() ([]Data.ReadActivity, error) {
	return []Data.ReadActivity{}, nil
}

func (s MockDBRepo) UpdateSitesAndArticles(sites []Data.WebSite, articles []Data.Article) error {
	return nil
}
