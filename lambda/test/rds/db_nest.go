package RDS

import (
	"errors"
	"log"

	"github.com/mmcdole/gofeed"
	"gorm.io/gorm"
)

// 入れ子に対して検索できるかテストする
func DbNestedStructTest(DB *gorm.DB) (bool, error) {
	// テーブル作成
	err := DB.Debug().AutoMigrate(&Site{}, &Feed{})
	if err != nil {
		log.Println("AutoMigrate Error:", err)
		return false, err
	}
	// GIGAZINEのRSSを取得する
	site, feeds, err := GetGIGAZINE()
	if err != nil {
		return false, err
	}
	// インサートする
	DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Debug().Create(&site)
		if result.Error != nil {
			log.Println("Insert Error:", result.Error)
			tx.Rollback()
			return result.Error
		}
		// result = tx.Create(&feeds)
		// if result.Error != nil {
		// 	log.Println("Insert Error:", result.Error)
		// 	tx.Rollback()
		// 	return result.Error
		// }
		// Commitに失敗した場合はロールバックされる
		return tx.Commit().Error
	})
	// そもそもインサート出来ているのか確認
	// レコード数を取得する
	var count int64
	DB.Model(&Site{}).Count(&count)
	var FeedCount int64
	DB.Model(&Feed{}).Count(&FeedCount)
	if count == 0 {
		log.Println("Insert Error: Not Insert")
		return false, errors.New("not insert")
	}
	log.Println("Insert Success Site Count:", count)
	log.Println("Insert Success Feed Count:", FeedCount)
	// ちゃんとインサートされている
	// 入れ子での検索/合格条件はTarget Titleがfeed[0]のタイトルと一致すること
	targetTile := feeds[0].Title
	// targetSiteTitle := "GIGAZINE"
	// 検索条件をログに出力
	log.Println("Target Feed Title:", targetTile)
	var res_feeds []Feed
	var res_site Site
	// 色々な書き方を試す
	// SQL文を見るとPreLoadが実行されていない
	// result := DB.Debug().Preload("SiteFeeds").Find(&res_site)
	// アソシエーションを試す
	result := DB.Debug().Preload("Feeds",&Feed{Title: targetTile}).Find(&res_site)
	if result.Error != nil {
		log.Println("Target Title:", targetTile)
		log.Println("Feed Match Error:", result.Error)
		// return false, result.Error
		return false, errors.New("not found")
	}
	// 何も取得できなかった場合はエラー
	if result.RowsAffected == 0 {
		log.Println("Target Title:", targetTile)
		log.Println("Feed Match Error: Not Found")
		return false, errors.New("not found")
	}
	// テーブルごと削除
	err = DB.Migrator().DropTable(&Site{})
	if err != nil {
		log.Println("Delete table error:", err)
	}
	err = DB.Migrator().DropTable(&Feed{})
	if err != nil {
		log.Println("Delete table error:", err)
	}
	// 検索結果が一致しない場合はエラー
	if res_feeds[0].Title != targetTile {
		log.Println("Target Title: ", targetTile)
		log.Println("Result Title: ", res_feeds[0].Title)
		log.Println("Feed Match Error: Not Match")
		return false, errors.New("not match")
	}
	return true, nil
}

func GetGIGAZINE() (Site, []Feed, error) {
	// GIGAZINEのURL
	url := "https://gigazine.net/news/rss_2.0/"
	// RSSを取得する
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return Site{}, nil, err
	}
	// RSS数が0の場合はエラー
	if len(feed.Items) == 0 {
		return Site{}, nil, errors.New("RSS is empty")
	}
	// Imageがない場合は空文字を入れる
	if feed.Image == nil {
		feed.Image = &gofeed.Image{URL: ""}
	}
	site := Site{
		SiteName:    "GIGAZINE",
		SiteUrl:     feed.Link,
		RssUrl:      url,
		IconUrl:     feed.Image.URL,
		Description: feed.Description,
	}
	// RSSをSiteFeed型の配列に変換する
	var siteFeeds []Feed
	for _, item := range feed.Items {
		// Imageがない場合は空文字を入れる
		if item.Image == nil {
			item.Image = &gofeed.Image{URL: ""}
		}
		siteFeeds = append(siteFeeds, Feed{
			Title:       item.Title,
			Description: item.Description,
			Url:         item.Link,
			IconUrl:     item.Image.URL,
			PublishedAt: *item.PublishedParsed,
			// Site:        &site,
		})
	}
	log.Println("site Title:", feed.Title)
	// RSSを含めたサイト情報を返す
	site.Feeds = siteFeeds
	return site, siteFeeds, nil
}
