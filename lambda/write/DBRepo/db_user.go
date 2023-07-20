package DBRepo

import (
	"read/Data"
	"time"

	"gorm.io/gorm"
)

// GORMで使う構造型のメンバは大文字で始めないといけない
type User struct {
	gorm.Model
	UserName         string
	UserUniqueID     int `gorm:"uniqueIndex"`
	AccountType      string
	Country          string
	ClientConfig     ClientConfig
	ApiActivity      []ApiActivity
	FavoriteSite     []FavoriteSite
	FavoriteArticle  []FavoriteArticle
	SubscriptionSite []SubscriptionSite
	ReadHistory      []ReadHistory
	SearchHistory    []SearchHistory
}

type ClientConfig struct {
	gorm.Model
	UserID                 uint
	ArticleRefreshInterval int
	UIConfig               UiConfig
	ApiConfig              ApiConfig
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

type FavoriteArticle struct {
	gorm.Model
	UserID    uint
	ArticleID uint
}

type SubscriptionSite struct {
	gorm.Model
	UserID      uint
	FolderIndex int
	FolderName  string
	SiteIndex   int
	SiteID      uint
}

type SearchHistory struct {
	gorm.Model
	UserID     uint
	SearchWord string
}

type ReadHistory struct {
	gorm.Model
	UserID         uint
	ActivityType   string
	ArticleID      uint
	SiteID         uint
	AccessAt       time.Time
	AccessPlatform string
	AccessIP       string
}

type ApiConfig struct {
	gorm.Model
	ClientConfigID              uint
	AccountType                 string
	FetchArticleRequestInterval int
	FetchArticleRequestLimit    int
	FetchTrendRequestInterval   int
	FetchTrendRequestLimit      int
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

// DB構造体からAPI構造体への変換
func ConvertToUserConfig(dbCfg User) (resUserCfg Data.UserConfig) {
	uiCfg := UiConfig{
		ClientConfigID:    dbCfg.ID,
		MobileTextSize:    dbCfg.ClientConfig.UIConfig.MobileTextSize,
		TabletTextSize:    dbCfg.ClientConfig.UIConfig.TabletTextSize,
		ThemeColor:        dbCfg.ClientConfig.UIConfig.ThemeColor,
		ThemeMode:         dbCfg.ClientConfig.UIConfig.ThemeMode,
		DrawerMenuOpacity: dbCfg.ClientConfig.UIConfig.DrawerMenuOpacity,
	}
	apiCfg := ApiConfig{
		ClientConfigID:              dbCfg.ID,
		AccountType:                 dbCfg.ClientConfig.ApiConfig.AccountType,
		FetchArticleRequestInterval: dbCfg.ClientConfig.ApiConfig.FetchArticleRequestInterval,
		FetchArticleRequestLimit:    dbCfg.ClientConfig.ApiConfig.FetchArticleRequestLimit,
		FetchTrendRequestInterval:   dbCfg.ClientConfig.ApiConfig.FetchTrendRequestInterval,
		FetchTrendRequestLimit:      dbCfg.ClientConfig.ApiConfig.FetchTrendRequestLimit,
	}
	// UniqueIDをintからstringに変換
	uniqueId := string(dbCfg.UserUniqueID)
	// 検索履歴をData.SearchHistoryの配列に変換
	var searchHistory []Data.SearchHistory
	for _, searchHistoryDb := range dbCfg.SearchHistory {
		searchHistory = append(searchHistory, Data.SearchHistory{
			SearchWord: searchHistoryDb.SearchWord,
			// 検索履歴の日付をRFC3339に変換
			CreatedAt: searchHistoryDb.CreatedAt.Format(time.RFC3339),
		})
	}
	// 読んだ記事をData.ReadHistoryの配列に変換
	var readHistory []Data.ReadHistory
	for _, readHistoryDb := range dbCfg.ReadHistory {
		readHistory = append(readHistory, Data.ReadHistory{
			ActivityType:   readHistoryDb.ActivityType,
			ArticleID:      readHistoryDb.ArticleID,
			SiteID:         readHistoryDb.SiteID,
			AccessAt:       readHistoryDb.AccessAt.Format(time.RFC3339),
			AccessPlatform: readHistoryDb.AccessPlatform,
			AccessIP:       readHistoryDb.AccessIP,
		})
	}
	// 購読サイトをData.SubscribeWebSiteの配列に変換
	var subscribeWebSite []Data.SubscribeWebSite
	for _, subscribeWebSiteDb := range dbCfg.SubscriptionSite {
		//これだとWebSiteを変換する為に都度クエリが走るから、ここはSiteIDだけ保持してクライアント側で必要になったらSite情報を取得するようにする
		subscribeWebSite = append(subscribeWebSite, Data.SubscribeWebSite{
			SiteID:      subscribeWebSiteDb.SiteID,
			FolderIndex: subscribeWebSiteDb.FolderIndex,
			FolderName:  subscribeWebSiteDb.FolderName,
			SiteIndex:   subscribeWebSiteDb.SiteIndex,
		})
	}
	// お気に入りサイトをData.FavoriteSiteの配列に変換
	var favoriteSite []Data.FavoriteSite
	for _, favoriteSiteDb := range dbCfg.FavoriteSite {
		favoriteSite = append(favoriteSite, Data.FavoriteSite{
			SiteID:    favoriteSiteDb.SiteID,
			CreatedAt: favoriteSiteDb.CreatedAt.Format(time.RFC3339),
		})
	}
	// お気に入り記事をData.FavoriteArticleの配列に変換
	var favoriteArticle []Data.FavoriteArticle
	for _, favoriteArticleDb := range dbCfg.FavoriteArticle {
		favoriteArticle = append(favoriteArticle, Data.FavoriteArticle{
			ArticleID: favoriteArticleDb.ArticleID,
			CreatedAt: favoriteArticleDb.CreatedAt.Format(time.RFC3339),
		})
	}
	// AppConfigをData.UserConfigに変換
	return Data.UserConfig{
		UserName:         dbCfg.UserName,
		UserUniqueID:     uniqueId,
		AccountType:      dbCfg.AccountType,
		SearchHistory:    searchHistory,
		SubscribeWebSite: subscribeWebSite,
		FavoriteSite:     favoriteSite,
		FavoriteArticle:  favoriteArticle,
		ReadHistory:      readHistory,
	}
}

func ConvertToDbUserConfig(apiCfg Data.UserConfig) (resDbUserCfg User) {
	//
	return User{}
}
