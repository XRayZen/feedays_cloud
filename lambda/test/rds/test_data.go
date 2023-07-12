package RDS

import (
	"time"

	"gorm.io/gorm"
)

// 構造体の入れ子に対して検索できるかテストする為に構造体をつくる
// サイトの中に記事リストがある
// フォーリンキーは特に設定せずライブラリに任せる
type Site struct {
	// gorm.Modelをつけると、idとCreatedAtとUpdatedAtとDeletedAtが作られる
	gorm.Model
	SiteName    string
	SiteUrl     string
	RssUrl      string
	IconUrl     string
	Description string
	Feeds       []Feed
}

type Feed struct {
	gorm.Model
	SiteID      uint
	Title       string `gorm:"column:title"`
	Url         string
	IconUrl     string
	Description string
	PublishedAt time.Time
}
