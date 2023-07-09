package Repo

import (
	"time"

	"gorm.io/gorm"
)

type Site struct {
	// gorm.Modelをつけると、idとCreatedAtとUpdatedAtとDeletedAtが作られる
	// gormでは、gorm.Modelでdeleted_atを作成している場合、Deleteすると、自動的に論理削除になるという仕様
	gorm.Model
	site_name        string
	site_url         string
	rss_url          string
	icon_url         string
	description      string
	site_feeds       []SiteFeed      `gorm:"foreignKey:site_id"`
	tags             []Tag           `gorm:"foreignKey:site_id"`
}

type SiteFeed struct {
	gorm.Model
	site_id      uint
	feed_index   int
	title        string
	url          string
	icon_url     string
	description  string
	published_at time.Time
}

type Tag struct {
	gorm.Model
	tag_name string
	site_id  uint
}

type ExploreCategory struct {
	gorm.Model
	category_name string
	description   string
	country       string
}

type SiteRanking struct {
	gorm.Model
	country             string
	explore_category_id uint
	site                Site `gorm:"foreignKey:site_id"`
	ranking_index       int
}

type FeedRanking struct {
	gorm.Model
	country             string
	explore_category_id uint
	site                Site     `gorm:"foreignKey:site_id"`
	feed                SiteFeed `gorm:"foreignKey:feed_id"`
	ranking_index       int
}
