package Repo

import (
	"batch/Data"
	"errors"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
)

func TestBatchDbRepo(t *testing.T) {
	dbRepo := InitDataBase()
	// Userを取得する
	var user User
	// Userが取得出来ない
	if err := DBMS.Where("user_unique_id = ?", "0000").Find(&user).Error; err != nil {
		t.Errorf("failed to get user")
	}
	// ReadHistoryをUserと紐づけて保存する
	if err := DBMS.Model(&user).Association("ReadHistories").Append(&ReadHistory{
		// GIGAZINEの記事を入れておく
		Link:     "https://gigazine.net/news/20201001-ai-robot-artist/",
		AccessAt: time.Now(),
	}); err != nil {
		t.Errorf("failed to append read history")
	}
	// GIGAZINEを入れておく
	site, articles, err := GetGIGAZINE()
	if err != nil {
		panic("failed to get GIGAZINE")
	}
	// API構造体からDB構造体に変換する
	dbSite := convertApiSiteToDb(site, articles)
	DBMS.Create(&dbSite)
	t.Run("Batch", func(t *testing.T) {
		// バッチ処理のDB操作をテストする
		sites, err := dbRepo.FetchAllSites()
		if err != nil || len(sites) == 0 {
			t.Errorf("failed to fetch all sites")
		}
		readHists, err := dbRepo.FetchAllReadHistories()
		if err != nil || len(readHists) == 0 {
			t.Errorf("failed to fetch all read histories")
		}
		// 記事を作る
		nowArticleTime := time.Now().UTC().Add(time.Hour).AddDate(0, 0, 10).Format(time.RFC3339)
		article := Data.Article{
			Link:        "https://gigazine.net/news/20201001-ai-robot-artist/",
			Title:       "TestArticle",
			Description: "AIが描いた絵がオークションで約1億円で落札される",
			PublishedAt: nowArticleTime,
		}
		sites[0].LastModified = nowArticleTime
		if err := dbRepo.UpdateSiteAndArticle(sites[0], []Data.Article{article}); err != nil {
			t.Errorf("failed to update site and article")
		}
		// 記事が更新されているか確認する
		updatedSite, err := dbRepo.SearchSiteLatestArticle(site.SiteURL, 100)
		if err != nil || updatedSite[0].PublishedAt != nowArticleTime {
			t.Errorf("failed to search site latest article")
		}
		// サイトの更新時間が更新されているか確認する
		// サイトを取得する
		var dbSite Site
		if err := DBMS.Where("site_url = ?", site.SiteURL).Find(&dbSite).Error; err != nil {
			t.Errorf("failed to get site")
		}
		// サイトの更新時間が更新されているか確認する
		timeLastModified, err := time.Parse(time.RFC3339, nowArticleTime)
		if err != nil {
			t.Errorf("failed to parse time")
		}
		if dbSite.LastModified != timeLastModified {
			t.Errorf("failed to update site last modified")
		}
	})
}

// テスト用のモックデータを生成する
func InitDataBase() DBRepository {
	dbRepo := DBRepoImpl{}
	// MockModeでRDSではなくインメモリーsqliteに接続する
	if err := dbRepo.ConnectDB(true); err != nil {
		panic("failed to connect database")
	}
	if err := dbRepo.AutoMigrate(); err != nil {
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
	return dbRepo
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