package Repo

import (
	"errors"
	"read/Data"
	"time"

	"testing"

	"github.com/mmcdole/gofeed"
)

// DBRepoImplをテストしてSQLクエリが正しく動作するか検証する

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

func TestRead(t *testing.T) {
	// テストを実行する
	t.Run("Read", func(t *testing.T) {
		dbRepo := InitDataBase()
		//カテゴリーが取得できるか検証する
		result, err := dbRepo.FetchExploreCategories("JP")
		if err != nil {
			t.Errorf("failed to fetch categories")
		}
		if len(result) == 0 {
			t.Errorf("failed to fetch categories")
		}
		// SQLiteはインメモリーなので、データが消えるから後始末は不要
	})
}




