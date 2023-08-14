package DBRepo

import (
	"errors"
	"testing"
	"time"
	"user/Data"

	"github.com/mmcdole/gofeed"
)

// テスト用のモックデータを生成する
func InitDataBase() DBRepo {
	dbRepo := DBRepoImpl{}
	// MockModeでRDSではなくインメモリーsqliteに接続する
	if err := dbRepo.ConnectDB(true); err != nil {
		panic("failed to connect database")
	}
	DBMS.AutoMigrate(&User{})
	DBMS.AutoMigrate(&FavoriteSite{})
	DBMS.AutoMigrate(&FavoriteArticle{})
	DBMS.AutoMigrate(&SubscriptionSite{})
	DBMS.AutoMigrate(&SearchHistory{})
	DBMS.AutoMigrate(&ReadHistory{})
	DBMS.AutoMigrate(&ApiLimitConfig{})
	DBMS.AutoMigrate(&UiConfig{})

	DBMS.AutoMigrate(&Site{})
	DBMS.AutoMigrate(&Article{})
	DBMS.AutoMigrate(&Tag{})
	DBMS.AutoMigrate(&ExploreCategory{})

	// カテゴリを生成する
	var categories = []ExploreCategory{
		{
			CategoryName: "CategoryName",
			Country:      "JP",
		},
	}
	// カテゴリを保存する
	DBMS.Create(&categories)
	return dbRepo
}

// GIGAZINEのサイトを取得して生成する
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

func TestDbRepoTest(t *testing.T) {
	// まずはUserを生成する
	dbRepo := InitDataBase()
	// 登録するUserを生成する
	user := Data.UserConfig{
		UserName:      "MockUser",
		UserUniqueID:  "0000",
		AccountType:   "Free",
		Country:       "JP",
		SearchHistory: []Data.SearchHistory{},
		ClientConfig:  Data.ClientConfig{},
		ReadHistory:   []Data.ReadHistory{},
	}
	site, articles, err := GetGIGAZINE()
	if err != nil {
		t.Errorf("failed to get site: %v", err)
	}
	dbSite := convertApiSiteToDb(site, articles)
	result := DBMS.Create(&dbSite)
	if result.Error != nil {
		t.Errorf("failed to create site: %v", result.Error)
	}
	// Userを取得する
	t.Run("GetUser", func(t *testing.T) {
		user.ClientConfig.UiConfig.DrawerMenuOpacity = 0.5
		user.ClientConfig.UiConfig.ThemeMode = "Dark"
		// Userを登録する
		if err := dbRepo.RegisterUser(user); err != nil {
			t.Errorf("failed to create user: %v", err)
		}
		// Userを取得する
		user, err := dbRepo.SearchUserConfig(user.UserUniqueID)
		// 取得出来たのがMockUserか確認する
		if err != nil || user.UserName != "MockUser" || user.ClientConfig.UiConfig.DrawerMenuOpacity != 0.5 || user.ClientConfig.UiConfig.ThemeMode != "Dark" {
			t.Errorf("failed to get user: %v", err)
		}
		// Userを更新する
		user.ClientConfig.UiConfig.DrawerMenuOpacity = 0.6
		if err := dbRepo.UpdateUser(user.UserUniqueID, user); err != nil {
			t.Errorf("failed to update user: %v", err)
		}
		// Userが更新されているか確認する
		user, err = dbRepo.SearchUserConfig(user.UserUniqueID)
		// 数字が反映されていない場合はエラー
		if err != nil || user.ClientConfig.UiConfig.DrawerMenuOpacity != 0.6 {
			// 検証に失敗した場合はエラー
			t.Errorf("failed to update validation: %v", err)
		}
		// 閲覧履歴を追加する
		readActivity := Data.ReadHistory{
			Link:           "link",
			AccessAt:       time.Now().Format(time.RFC3339),
			AccessPlatform: "PC",
			AccessIP:       "10.9.9.9",
		}
		if err := dbRepo.AddReadHistory(user.UserUniqueID, readActivity); err != nil {
			t.Errorf("failed to add read history: %v", err)
		}
		//閲覧履歴を取得して検証する
		user, err = dbRepo.SearchUserConfig(user.UserUniqueID)
		if err != nil || len(user.ReadHistory) == 0 {
			t.Errorf("failed to get read history: %v", err)
		}
		// 検索履歴を追加する
		words, err := dbRepo.ModifySearchHistory(user.UserUniqueID, "SearchWord", true)
		if err != nil || len(words) == 0 {
			t.Errorf("failed to add search history: %v", err)
		}
		// 検索履歴を削除する
		words, err = dbRepo.ModifySearchHistory(user.UserUniqueID, "SearchWord", false)
		if err != nil || len(words) != 0 {
			t.Errorf("failed to delete search history: %v", err)
		}
		// お気に入りサイトを追加する
		if err := dbRepo.ModifyFavoriteSite(user.UserUniqueID, site.SiteURL, true); err != nil {
			t.Errorf("failed to add favorite site: %v", err)
		}
		// お気に入りサイトを取得する
		user, err = dbRepo.SearchUserConfig(user.UserUniqueID)
		if err != nil || len(user.FavoriteSite) == 0 {
			t.Errorf("failed to get favorite site: %v", err)
		}
		// お気に入りサイトを削除する
		if err := dbRepo.ModifyFavoriteSite(user.UserUniqueID, site.SiteURL, false); err != nil {
			t.Errorf("failed to delete favorite site: %v", err)
		}
		// お気に入り記事を追加する
		if err := dbRepo.ModifyFavoriteArticle(user.UserUniqueID, articles[0].Link, true); err != nil {
			t.Errorf("failed to add favorite article: %v", err)
		}
		// お気に入り記事を取得する
		user, err = dbRepo.SearchUserConfig(user.UserUniqueID)
		if err != nil || len(user.FavoriteArticle) == 0 {
			t.Errorf("failed to get favorite article: %v", err)
		}
		// お気に入り記事を削除する
		if err := dbRepo.ModifyFavoriteArticle(user.UserUniqueID, articles[0].Link, false); err != nil {
			t.Errorf("failed to delete favorite article: %v", err)
		}
		// API設定を追加する
		if err := dbRepo.ModifyApiRequestLimit("Add", Data.ApiConfig{
			AccountType:                 "Free",
			FetchArticleRequestInterval: 1000,
			FetchTrendRequestInterval:   2000,
		}); err != nil {
			t.Errorf("failed to add api config: %v", err)
		}
		// API設定を取得する
		apiConfig, err := dbRepo.FetchAPIRequestLimit(user.UserUniqueID)
		if err != nil || apiConfig.FetchArticleRequestInterval != 1000 || apiConfig.FetchTrendRequestInterval != 2000 {
			t.Errorf("failed to get api config: %v", err)
		}
	})
}
