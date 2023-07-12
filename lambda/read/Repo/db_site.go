package Repo

import (
	"time"

	"gorm.io/gorm"
)

type Site struct {
	// gorm.Modelをつけると、idとCreatedAtとUpdatedAtとDeletedAtが作られる
	// gormでは、gorm.Modelでdeleted_atを作成している場合、Deleteすると、自動的に論理削除になるという仕様
	gorm.Model
	SiteName    string
	SiteUrl     string
	RssUrl      string
	IconUrl     string
	Description string
	SiteFeeds   []SiteFeed
	Tags        []Tag
	Category    string
}

type SiteFeed struct {
	gorm.Model
	SiteID      uint
	FeedIndex   int
	Title       string
	Url         string
	IconUrl     string
	Description string
	PublishedAt time.Time
}

type Tag struct {
	gorm.Model
	TagName string
	SiteID  uint
}

type ExploreCategory struct {
	gorm.Model
	CategoryName string
	Description  string
	Country      string
}

type SiteRanking struct {
	gorm.Model
	Country           string
	ExploreCategoryID uint
	SiteID            uint
	RankingIndex      int
}

type FeedRanking struct {
	gorm.Model
	Country           string
	ExploreCategoryID uint
	SiteID            uint
	SiteFeedID        uint
	RankingIndex      int
}
