package Repo

import (
	// "log"
	"site/Data"
	"time"
)

type MockDBRepo struct {
	// articles []Data.Article
}

// DB接続
func (s MockDBRepo) ConnectDB(isMock bool) error {
	return nil
}

// DBオートマイグレート
func (s MockDBRepo) AutoMigrate() error {
	return nil
}

// readで使う
func (s MockDBRepo) SearchUserConfig(user_unique_Id string, isPreloadRelatedTables bool) (Data.UserConfig, error) {
	return Data.UserConfig{
		UserName:     "UserName",
		UserUniqueID: "UserUniqueID",
	}, nil
}

func (s MockDBRepo) FetchExploreCategories(country string) (resExp []Data.ExploreCategory, err error) {
	return []Data.ExploreCategory{
		{
			CategoryName: "CategoryName",
		},
	}, nil
}

// モック用変数 テストする時に変更する
var MockSiteLastModified int

func (s MockDBRepo) UpdateSiteAndArticle(site Data.WebSite, articles []Data.Article) error {
	return nil
}

func (s MockDBRepo) SearchReadActivityByTime(from time.Time, to time.Time) ([]Data.ReadHistory, error) {
	return nil, nil
}

// ランキングを更新
func (s MockDBRepo) UpdateRanking() error {
	return nil
}

// IsExistSite(site_url string) bool
// サイトが存在するかどうか
func (s MockDBRepo) IsExistSite(site_url string) bool {
	switch site_url {
	case "https://automaton-media.com/":
		return true
	case "https://gigazine.net/":
		return true
	}
	return false
}

// RegisterSite(site Data.WebSite, articles []Data.Article) error
// サイトを登録
func (s MockDBRepo) RegisterSite(site Data.WebSite, articles []Data.Article) error {
	return nil
}

// SearchArticlesByKeyword(keyword string) ([]Data.Article, error)
// キーワードで検索
func (s MockDBRepo) SearchArticlesByKeyword(keyword string) ([]Data.Article, error) {
	if keyword == "Found" {
		return []Data.Article{
			{
				Title:       "Found",
				Link:        "https://example.com",
				Site:        "https://example.com",
				PublishedAt: "2021-01-01T00:00:00+09:00",
			},
		}, nil
	}
	return nil, nil
}

// SearchSiteLatestArticle(site_url string, get_count int) ([]Data.Article, error)
// サイトの最新記事を取得
func (s MockDBRepo) SearchSiteLatestArticle(site_url string, get_count int) ([]Data.Article, error) {
	return []Data.Article{
		{
			Title:       "Found",
			Link:        "https://example.com",
			Site:        "https://example.com",
			PublishedAt: "2021-01-01T00:00:00+09:00",
		},
	}, nil
}

// モック用変数
var mockArticles = []Data.Article{}

// SearchArticlesByTimeAndOrder(siteUrl string, lastModified time.Time, get_count int, isNew bool) ([]Data.Article, error)
// 時間と順番で検索
func (s MockDBRepo) SearchArticlesByTimeAndOrder(siteUrl string, lastModified time.Time, get_count int, isNew bool) ([]Data.Article, error) {
	// 更新日時より新しい記事を返す
	var articles []Data.Article
	for _, article := range mockArticles {
		articleTime, _ := time.Parse(time.RFC3339, article.PublishedAt)
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

// SearchSiteByName(siteName string) ([]Data.WebSite, error)
// サイト名で検索
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

func (s MockDBRepo) IsSubscribeSite(user_id string, site_url string) bool {
	return true
}

// SearchSiteByURL(siteURL string) ([]Data.WebSite, error)
// サイトURLで検索
func (s MockDBRepo) SearchSiteByUrl(siteURL string) (Data.WebSite, error) {
	// テスト用にダミーのWebSiteを返す
	switch siteURL {
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

// SearchSiteByCategory(category string) ([]Data.WebSite, error)
// カテゴリーで検索
func (s MockDBRepo) SearchSiteByCategory(category string) ([]Data.WebSite, error) {
	return nil, nil
}

// SubscribeSite(user_unique_id string, site_url string, is_subscribe bool) error
// サイトを購読
func (s MockDBRepo) SubscribeSite(user_unique_id string, site_url string, is_subscribe bool) error {
	return nil
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

// ChangeSiteCategory(user_unique_id string, site_url string, category_name string) error
func (s MockDBRepo) ChangeSiteCategory(user_unique_id string, site_url string, category_name string) error {
	return nil
}

// 	DeleteSite(site_url string) error
func (s MockDBRepo) DeleteSite(site_url string) error {
	return nil
}

// DeleteSiteByUnscoped(site_url string) error
func (s MockDBRepo) DeleteSiteByUnscoped(site_url string) error {
	return nil
}

// ModifyExploreCategory(category Data.ExploreCategory, is_add_or_remove bool) error
func (s MockDBRepo) ModifyExploreCategory(modify_type string, category Data.ExploreCategory) error{
	return nil
}
