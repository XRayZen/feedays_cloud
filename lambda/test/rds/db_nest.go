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
	DB.AutoMigrate(&DbTestSite{})
	DB.AutoMigrate(&DbTestSiteFeed{})
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
	result := DB.Where(&DbTestSite{site_name: targetSiteTitle}).Where(&DbTestSiteFeed{title: targetTile}).First(&siteFeed)
	if result.Error != nil {
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
		return false, errors.New("Not match")
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
	// RSSをSiteFeed型の配列に変換する
	var siteFeeds []DbTestSiteFeed
	index := 0
	for _, item := range feed.Items {
		// Imageがない場合は空文字を入れる
		if item.Image == nil {
			item.Image = &gofeed.Image{URL: ""}
		}
		siteFeeds = append(siteFeeds, DbTestSiteFeed{
			title:        item.Title,
			feed_index:   index,
			description:  item.Description,
			url:          item.Link,
			icon_url:   item.Image.URL,
			published_at: *item.PublishedParsed,
		})
		index++
	}
	log.Println("site Title:", feed.Title)
	// Imageがない場合は空文字を入れる
	if feed.Image == nil {
		feed.Image = &gofeed.Image{URL: ""}
	}
	// RSSを含めたサイト情報を返す
	return DbTestSite{
		site_name:   feed.Title,
		site_url:    feed.Link,
		rss_url:     url,
		icon_url:    feed.Image.URL,
		description: feed.Description,
		site_feeds:  siteFeeds,
	}, siteFeeds, nil
}
