package Repo

import (
	// "log"
	"batch/Data"
	"time"
)

type MockDBRepo struct {
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
		UserName:   "UserName",
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

// SearchSiteLatestArticle(site_url string, get_count int) ([]Data.Article, error)
func (s MockDBRepo) SearchSiteLatestArticle(site_url string, get_count int) ([]Data.Article, error) {
	return nil, nil
}
