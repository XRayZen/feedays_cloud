package Repo

import "gorm.io/gorm"

// GORMで使う構造型のメンバは大文字で始めないといけない
type User struct {
	gorm.Model
	UserName         string
	UserUniqueID     int `gorm:"uniqueIndex"`
	AccountType      string
	ApiActivity      ApiActivity
	FavoriteSite     []FavoriteSite
	SubscriptionSite []SubscriptionSite
	ReadHistory      []ReadHistory
	SearchHistory    []SearchHistory
	ClientConfig     ClientConfig
}

type ClientConfig struct {
	gorm.Model
	UserID              uint
	FeedRefreshInterval int
	UIConfig            UiConfig
	ApiConfig           ApiConfig
}

type ApiActivity struct {
	gorm.Model
	UserID         uint
	ActivityType   string
	ActivityTime   string
	AccessPlatform string
	AccessIP       string
}

type FavoriteSite struct {
	gorm.Model
	UserID uint
	SiteID uint
}

type SubscriptionSite struct {
	gorm.Model
	UserID   uint
	site_url string
}

type SearchHistory struct {
	gorm.Model
	UserID     uint
	SearchWord string
}

type ReadHistory struct {
	gorm.Model
	UserID  uint
	SiteURL string
	FeedURL string
}

type ApiConfig struct {
	gorm.Model
	ClientConfigID            uint
	AccountType               string
	FetchFeedRequestInterval  int
	FetchFeedRequestLimit     int
	FetchTrendRequestInterval int
	FetchTrendRequestLimit    int
}

type UiConfig struct {
	gorm.Model
	ClientConfigID    uint
	MobileTextSize    int
	TabletTextSize    int
	ThemeColor        string
	ThemeMode         string
	DrawerMenuOpacity float32
}
