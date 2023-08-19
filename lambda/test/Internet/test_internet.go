package Internet

import (
	"errors"
	"github.com/mmcdole/gofeed"
)

// GIGAZINEのRSSを取得する
func GetGIGAZINE() (string, error) {
	// GIGAZINEのURL
	url := "https://gigazine.net/news/rss_2.0/"
	// RSSを取得する
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return "", err
	}
	// RSS数が0の場合はエラー
	if len(feed.Items) == 0 {
		return "", errors.New("RSS is empty")
	}
	// RSSのタイトルを返す
	return feed.Items[0].Title, nil
}
