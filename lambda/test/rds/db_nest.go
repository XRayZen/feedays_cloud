package RDS

import (
	"errors"
	"log"

	"github.com/mmcdole/gofeed"
	"gorm.io/gorm"
)

// 入れ子に対して検索できるかテストする
func DbNestedStructTest() (bool, error) {
	DB, err := GormConnect()
	if err != nil {
		return false, err
	}
	// テーブル作成
	err = DB.AutoMigrate(&DbTestSite{}, &DbTestSiteFeed{})
	if err != nil {
		log.Println("AutoMigrate Error:", err)
		return false, err
	}
	// GIGAZINEのRSSを取得する
	site, feeds, err := GetGIGAZINE()
	if err != nil {
		return false, err
	}
	targetSiteTitle := site.site_name
	// トランザクション開始
	DB.Transaction(func(tx *gorm.DB) error {
		// トランザクション内でのデータベース処理を行う(ここでは `db` ではなく `tx` を利用する)
		if err := tx.Create(&site).Error; err != nil {
			// エラーが発生した場合はロールバックする
			tx.Rollback()
			return err
		}
		// エラーがなければコミットする
		return tx.Commit().Error
	})
	// 入れ子での検索
	targetTile := feeds[8].title
	var siteFeed DbTestSiteFeed
	// 色々な書き方を試す
	result := DB.Where(&DbTestSite{site_name: targetSiteTitle})
	if result.Error != nil {
		log.Println("Site Match Error:", result.Error)
		return false, result.Error
	}
	// サイト名が一致したが、フィードが一致しないのであればDBに入れる時にサイトは入れられたがフィードは入れられていない
	// プリロードを試してみる
	result = DB.Where(&DbTestSite{site_name: targetSiteTitle}).Preload("DbTestSiteFeeds").Where(&DbTestSiteFeed{
		title: targetTile,
	}).Find(&siteFeed)
	if result.Error != nil {
		log.Println("Target Title:", targetTile)
		log.Println("Feed Match Error:", result.Error)
		return false, result.Error
	}
	// テーブルごと削除
	err = DB.Migrator().DropTable(&DbTestSite{})
	if err != nil {
		log.Println("Delete table error:", err)
	}
	err = DB.Migrator().DropTable(&DbTestSiteFeed{})
	if err != nil {
		log.Println("Delete table error:", err)
	}
	// 検索結果が一致しない場合はエラー
	if siteFeed.title != targetTile {
		return false, errors.New("not match")
	}
	return true, nil
}

func GetGIGAZINE() (DbTestSite, []DbTestSiteFeed, error) {
	// GIGAZINEのURL
	url := "https://gigazine.net/news/rss_2.0/"
	// RSSを取得する
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return DbTestSite{}, nil, err
	}
	// RSS数が0の場合はエラー
	if len(feed.Items) == 0 {
		return DbTestSite{}, nil, errors.New("RSS is empty")
	}
	// Imageがない場合は空文字を入れる
	if feed.Image == nil {
		feed.Image = &gofeed.Image{URL: ""}
	}
	site := DbTestSite{
		site_name:   feed.Title,
		site_url:    feed.Link,
		rss_url:     url,
		icon_url:    feed.Image.URL,
		description: feed.Description,
	}
	// RSSをSiteFeed型の配列に変換する
	var siteFeeds []DbTestSiteFeed
	for _, item := range feed.Items {
		// Imageがない場合は空文字を入れる
		if item.Image == nil {
			item.Image = &gofeed.Image{URL: ""}
		}
		siteFeeds = append(siteFeeds, DbTestSiteFeed{
			title:        item.Title,
			description:  item.Description,
			url:          item.Link,
			icon_url:     item.Image.URL,
			published_at: *item.PublishedParsed,
			site:         site,
		})
	}
	log.Println("site Title:", feed.Title)
	// RSSを含めたサイト情報を返す
	site.site_feeds = siteFeeds
	return site, siteFeeds, nil
}
