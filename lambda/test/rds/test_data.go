package RDS

import (
	"time"

	"gorm.io/gorm"
)

// 構造体の入れ子に対して検索できるかテストする為に構造体をつくる
// サイトの中に記事リストがある
// フォーリンキーは特に設定せずライブラリに任せる
type DbTestSite struct {
	// gorm.Modelをつけると、idとCreatedAtとUpdatedAtとDeletedAtが作られる
	gorm.Model
	site_name   string
	site_url    string
	rss_url     string
	icon_url    string
	description string
	site_feeds  []DbTestSiteFeed
}

type DbTestSiteFeed struct {
	gorm.Model
	site         DbTestSite `gorm:"foreignKey:site_id"`
	title        string
	url          string
	icon_url     string
	description  string
	published_at time.Time
}
