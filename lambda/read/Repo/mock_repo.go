package Repo

import (
	// "log"
	"read/Data"
	"time"
)

type MockDBRepo struct {
	articles []Data.Article
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
			CategoryName: "Test",
		},
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

func (s MockDBRepo) SearchSiteByUrl(site_url string) (Data.WebSite, error) {
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
				Title:       "Found",
				Link:        "https://example.com",
				Site:        "https://example.com",
				PublishedAt: "2021-01-01T00:00:00+09:00",
			},
		}, nil
	}
	return nil, nil
}

func (s MockDBRepo) SearchArticlesByTime(siteUrl string, lastModified time.Time) ([]Data.Article, error) {
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

func (s MockDBRepo) SearchArticlesByTimeAndOrder(siteUrl string, lastModified time.Time, get_count int, isOlder bool) ([]Data.Article, error) {
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

// モック用変数 テストする時に変更する
var MockSiteLastModified int

// バッチ処理用
func (s MockDBRepo) FetchAllSites() ([]Data.WebSite, error) {
	// 今より30分前の日時を返す
	lastModifiedTime := time.Now().Add(-time.Minute * time.Duration(MockSiteLastModified))
	// それをRFC3339形式に変換
	lastModified := lastModifiedTime.Format(time.RFC3339)
	return []Data.WebSite{
		{
			SiteURL:      "https://automaton-media.com/",
			SiteRssURL:   "https://automaton-media.com/feed/",
			SiteName:     "AUTOMATON",
			LastModified: lastModified,
		},
		{
			SiteURL:      "https://gigazine.net/",
			SiteRssURL:   "https://gigazine.net/news/rss_2.0/",
			SiteName:     "GIGAZINE",
			LastModified: lastModified,
		},
	}, nil
}

func (s MockDBRepo) FetchAllReadHistories() ([]ReadHistory, error) {
	// モック用のReadActivityを生成して返す
	// 一番読まれたのはGIGAZINEの記事（架空）
	// 二番目に読まれたのはAUTOMATONの記事（架空）
	var readActivities []ReadHistory
	// GIGAZINEの記事を100回読んだことにする
	for i := 0; i < 100; i++ {
		ra := ReadHistory{
			// UserUniqueID: "Mock User",
			Link: "https://gigazine.net/article/20210101-mock-article/",
		}
		readActivities = append(readActivities, ra)
	}
	// AUTOMATONの記事を50回読んだことにする
	for i := 0; i < 50; i++ {
		ra := ReadHistory{
			// UserUniqueID: "Mock User",
			Link: "https://automaton-media.com/article/20210101-mock-article/",
		}
		readActivities = append(readActivities, ra)
	}

	return readActivities, nil
}

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
