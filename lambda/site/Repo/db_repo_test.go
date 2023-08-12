package Repo

import (
	"errors"
	"site/Data"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
)

// テスト用のモックデータを生成する
func InitDataBase() DBRepository {
	db_repo := DBRepoImpl{}
	// MockModeでRDSではなくインメモリーsqliteに接続する
	if err := db_repo.ConnectDB(true); err != nil {
		panic("failed to connect database")
	}
	if err := db_repo.AutoMigrate(); err != nil {
		panic("failed to migrate database")
	}
	// カテゴリを生成する
	var categories = []ExploreCategory{
		{
			CategoryName: "CategoryName",
			Country:      "JP",
		},
	}
	// カテゴリを保存する
	DBMS.Create(&categories)
	// Userを生成する
	var users = []User{
		{
			UserName:     "MockUser",
			UserUniqueID: "0000",
			Country:      "JP",
		},
	}
	// Userを保存する
	DBMS.Create(&users)
	return db_repo
}

// GIGAZINEのRSSを取得する
func GetGIGAZINE() (Data.WebSite, []Data.Article, error) {
	// GIGAZINEのURL
	url := "https://gigazine.net/news/rss_2.0/"
	// RSSを取得する
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return Data.WebSite{}, nil, err
	}
	// RSS数が0の場合はエラー
	if len(feed.Items) == 0 {
		return Data.WebSite{}, nil, errors.New("RSS is empty")
	}
	// Feedを記事に変換する
	articles := []Data.Article{}
	for _, v := range feed.Items {
		// Feedのカテゴリはタグにしておく
		category := ""
		if len(v.Categories) > 0 {
			category = v.Categories[0]
		}
		article := Data.Article{
			Title:       v.Title,
			Link:        v.Link,
			Description: v.Description,
			Category:    category,
			Site:        feed.Title,
			PublishedAt: v.PublishedParsed.UTC().Format(time.RFC3339),
		}
		articles = append(articles, article)
	}
	return Data.WebSite{
		SiteName:        feed.Title,
		SiteURL:         feed.Link,
		SiteRssURL:      url,
		SiteDescription: feed.Description,
		SiteTags:        feed.Categories,
		LastModified:    feed.UpdatedParsed.UTC().Format(time.RFC3339),
	}, articles, nil
}

// Heavy系のテスト
// サイトの登録・存在確認・取得
func TestDBQuerySite(t *testing.T) {
	db_repo := InitDataBase()
	// サイトの登録・存在確認・取得
	// サイト記事系の処理は
	t.Run("Write", func(t *testing.T) {
		site, articles, err := GetGIGAZINE()
		if err != nil {
			t.Errorf("failed to get GIGAZINE")
		}
		err = db_repo.RegisterSite(site, articles)
		if err != nil {
			t.Errorf("failed to register site : %v", err)
		}
		result := db_repo.IsExistSite(site.SiteURL)
		if !result {
			t.Errorf("failed to check site")
		}
		// サイト購読・確認
		err = db_repo.SubscribeSite("0000", site.SiteURL, true)
		if err != nil {
			t.Errorf("failed to subscribe site")
		}
		result = db_repo.IsSubscribeSite("0000", site.SiteURL)
		if !result {
			t.Errorf("failed to check subscribe site")
		}
		// サイトの最終更新日時が更新されているかテストする
		last_modified, err := db_repo.FetchSiteLastModified(site.SiteURL)
		if err != nil {
			t.Errorf("failed to fetch site last modified : %v", err)
		}
		// 時間が変換されていく中でローカル時間とずれるからUTCに変換して解決
		site_last_modified, err := time.Parse(time.RFC3339, site.LastModified)
		if err != nil || last_modified != site_last_modified {
			t.Errorf("failed to fetch site last modified")
		}
		// サイトURLをキーにサイトを検索
		result_site, err := db_repo.SearchSiteByUrl(site.SiteURL)
		if err != nil || result_site.SiteURL != site.SiteURL {
			t.Errorf("failed to search site by url")
		}
		// サイト名をキーにサイトを検索
		result_sites, err := db_repo.SearchSiteByName(site.SiteName)
		if err != nil || len(result_sites) == 0 {
			t.Errorf("failed to search site by name")
		}
		// 記事系のテストを行う
		// 記事の検索・存在確認・取得
		result_articles, err := db_repo.SearchArticlesByKeyword("ニュース")
		if err != nil && len(result_articles) == 0 {
			t.Errorf("failed to search articles by keyword")
		}
		result_articles, err = db_repo.SearchSiteLatestArticle(site.SiteURL, 100)
		// ちゃんと記事最新なのか確認する
		if err != nil || len(result_articles) == 0 || result_articles[0].PublishedAt != articles[0].PublishedAt {
			t.Errorf("failed to search site latest article")
		}
		// 今から五時間前から最新の記事を取得
		previous_time := time.Now().UTC().Add(-10 * time.Hour)
		result_articles, err = db_repo.SearchArticlesByTimeAndOrder(site.SiteURL, previous_time, 100, true)
		if err != nil && len(result_articles) == 0 {
			t.Errorf("failed to search articles by time and order")
		}
		// 正解を判定する
		// エラー条件は5時間前以上の記事があったらエラー
		for _, article := range result_articles {
			article_time, err := time.Parse(time.RFC3339, article.PublishedAt)
			if err != nil {
				t.Errorf("failed to parse article time")
			}
			if article_time.Before(time.Now().Add(-10 * time.Hour)) {
				t.Errorf("failed to search articles by time and order")
			}
		}
	})
}

// 記事系をテストする
func TestDBQueryArticle(t *testing.T) {
	db_repo := InitDataBase()
	// 記事の登録・存在確認・取得
	t.Run("Write", func(t *testing.T) {
		site, articles, err := GetGIGAZINE()
		if err != nil {
			t.Errorf("failed to get GIGAZINE")
		}
		err = db_repo.RegisterSite(site, articles)
		if err != nil {
			t.Errorf("failed to register site")
		}
		// 記事の検索・存在確認・取得
		result_articles, err := db_repo.SearchArticlesByKeyword("ニュース")
		if err != nil && len(result_articles) == 0 {
			t.Errorf("failed to search articles by keyword")
		}
		result_articles, err = db_repo.SearchSiteLatestArticle(site.SiteURL, 100)
		if err != nil && len(result_articles) == 0 {
			t.Errorf("failed to search site latest article")
		}
		// 今から五時間前から最新の記事を取得
		previous_time := time.Now().UTC().Add(-10 * time.Hour)
		result_articles, err = db_repo.SearchArticlesByTimeAndOrder(site.SiteURL, previous_time, 100, true)
		if err != nil && len(result_articles) == 0 {
			t.Errorf("failed to search articles by time and order")
		}
		// 正解を判定する
		// エラー条件は5時間前以上の記事があったらエラー
		for _, article := range result_articles {
			article_time, err := time.Parse(time.RFC3339, article.PublishedAt)
			if err != nil {
				t.Errorf("failed to parse article time")
			}
			if article_time.Before(time.Now().Add(-10 * time.Hour)) {
				t.Errorf("failed to search articles by time and order")
			}
		}
		// 今から10時間前よりも古い記事を取得
		result_articles, err = db_repo.SearchArticlesByTimeAndOrder(site.SiteURL, previous_time, 100, false)
		if err != nil && len(result_articles) == 0 {
			t.Errorf("failed to search articles by time and order")
		}
		// 正解を判定する
		// エラー条件は5時間前以下の記事があったらエラー
		for _, article := range result_articles {
			article_time, err := time.Parse(time.RFC3339, article.PublishedAt)
			if err != nil {
				t.Errorf("failed to parse article time")
			}
			if article_time.After(time.Now().Add(-10 * time.Hour)) {
				t.Errorf("failed to search articles by time and order")
			}
		}
	})
}
