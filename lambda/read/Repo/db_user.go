package Repo

import (
	"read/Data"
	"time"

	"gorm.io/gorm"
)

// GORMで使う構造型のメンバは大文字で始めないといけない
type User struct {
	gorm.Model
	UserName         string
	UserUniqueID     string `gorm:"uniqueIndex"`
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
	UserID    uint
	UiConfig  UiConfig
	ApiConfig ApiConfig
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
	searchAt   time.Time
}

type ReadHistory struct {
	gorm.Model
	UserID         uint
	ArticleID      uint
	SiteID         uint
	ActivityType   string
	AccessAt       time.Time
	AccessPlatform string
	AccessIP       string
}

type ApiConfig struct {
	gorm.Model
	ClientConfigID              uint
	RefreshArticleInterval      int
	FetchArticleRequestInterval int
	FetchArticleRequestLimit    int
	FetchTrendRequestInterval   int
	FetchTrendRequestLimit      int
}

type UiConfig struct {
	gorm.Model
	ClientConfigID        uint
	MobileTextSize        int
	TabletTextSize        int
	ThemeColorValue       int
	ThemeMode             string
	DrawerMenuOpacity     float64
	ArticleListFontSize   UiResponsiveFontSize
	ArticleDetailFontSize UiResponsiveFontSize
}

type UiResponsiveFontSize struct {
	gorm.Model
	UiConfigID uint
	Mobile     float64 `json:"mobile"`
	Tablet     float64 `json:"tablet"`
	Default    float64 `json:"defaultSize"`
}

// DB構造体からAPI構造体への変換
func ConvertToApiUserConfig(dbCfg User) (resUserCfg Data.UserConfig) {
	// UiResponsiveFontSizeをData.UiResponsiveFontSizeに変換
	articleListFontSize := Data.UiResponsiveFontSize{
		Mobile:  dbCfg.ClientConfig.UiConfig.ArticleListFontSize.Mobile,
		Tablet:  dbCfg.ClientConfig.UiConfig.ArticleListFontSize.Tablet,
		Default: dbCfg.ClientConfig.UiConfig.ArticleListFontSize.Default,
	}
	articleDetailFontSize := Data.UiResponsiveFontSize{
		Mobile:  dbCfg.ClientConfig.UiConfig.ArticleDetailFontSize.Mobile,
		Tablet:  dbCfg.ClientConfig.UiConfig.ArticleDetailFontSize.Tablet,
		Default: dbCfg.ClientConfig.UiConfig.ArticleDetailFontSize.Default,
	}
	uiCfg := Data.UiConfig{
		ThemeColorValue:       dbCfg.ClientConfig.UiConfig.ThemeColorValue,
		ThemeMode:             dbCfg.ClientConfig.UiConfig.ThemeMode,
		DrawerMenuOpacity:     dbCfg.ClientConfig.UiConfig.DrawerMenuOpacity,
		ArticleListFontSize:   articleListFontSize,
		ArticleDetailFontSize: articleDetailFontSize,
	}
	apiCfg := Data.ApiConfig{
		RefreshArticleInterval:      dbCfg.ClientConfig.ApiConfig.RefreshArticleInterval,
		FetchArticleRequestInterval: dbCfg.ClientConfig.ApiConfig.FetchArticleRequestInterval,
		FetchArticleRequestLimit:    dbCfg.ClientConfig.ApiConfig.FetchArticleRequestLimit,
		FetchTrendRequestInterval:   dbCfg.ClientConfig.ApiConfig.FetchTrendRequestInterval,
		FetchTrendRequestLimit:      dbCfg.ClientConfig.ApiConfig.FetchTrendRequestLimit,
	}
	// ClientConfigをData.ClientConfigに変換
	clientConfig := Data.ClientConfig{
		UiConfig:  uiCfg,
		ApiConfig: apiCfg,
	}
	// 検索履歴をData.SearchHistoryの配列に変換
	var searchHistory []Data.SearchHistory
	for _, searchHistoryDb := range dbCfg.SearchHistory {
		searchHistory = append(searchHistory, Data.SearchHistory{
			SearchWord: searchHistoryDb.SearchWord,
			// 検索履歴の日付をRFC3339に変換
			SearchAt: searchHistoryDb.CreatedAt.Format(time.RFC3339),
		})
	}
	// 読んだ記事をData.ReadHistoryの配列に変換
	var readHistory []Data.ReadHistory
	for _, readHistoryDb := range dbCfg.ReadHistory {
		readHistory = append(readHistory, Data.ReadHistory{
			UserID:         dbCfg.UserUniqueID,
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
	return Data.UserConfig{
		UserName:         dbCfg.UserName,
		UserUniqueID:     dbCfg.UserUniqueID,
		ClientConfig:     clientConfig,
		AccountType:      dbCfg.AccountType,
		SearchHistory:    searchHistory,
		SubscribeWebSite: subscribeWebSite,
		FavoriteSite:     favoriteSite,
		FavoriteArticle:  favoriteArticle,
		ReadHistory:      readHistory,
	}
}

// 検索履歴と閲覧記事履歴と購読サイトとお気に入りサイト・記事は変換しない
// 使用用途的には、ユーザー・API設定の更新時に使用するから
func ConvertToDbUserConfig(apiCfg Data.UserConfig) (resDbUserCfg User) {
	//フォントサイズをData.UiResponsiveFontSizeからDb.UiResponsiveFontSizeに変換
	articleListFontSize := UiResponsiveFontSize{
		Mobile:  apiCfg.ClientConfig.UiConfig.ArticleListFontSize.Mobile,
		Tablet:  apiCfg.ClientConfig.UiConfig.ArticleListFontSize.Tablet,
		Default: apiCfg.ClientConfig.UiConfig.ArticleListFontSize.Default,
	}
	articleDetailFontSize := UiResponsiveFontSize{
		Mobile:  apiCfg.ClientConfig.UiConfig.ArticleDetailFontSize.Mobile,
		Tablet:  apiCfg.ClientConfig.UiConfig.ArticleDetailFontSize.Tablet,
		Default: apiCfg.ClientConfig.UiConfig.ArticleDetailFontSize.Default,
	}
	uiCfg := UiConfig{
		ThemeColorValue:       apiCfg.ClientConfig.UiConfig.ThemeColorValue,
		ThemeMode:             apiCfg.ClientConfig.UiConfig.ThemeMode,
		DrawerMenuOpacity:     apiCfg.ClientConfig.UiConfig.DrawerMenuOpacity,
		ArticleListFontSize:   articleListFontSize,
		ArticleDetailFontSize: articleDetailFontSize,
	}
	apiCfgDb := ApiConfig{
		RefreshArticleInterval:      apiCfg.ClientConfig.ApiConfig.RefreshArticleInterval,
		FetchArticleRequestInterval: apiCfg.ClientConfig.ApiConfig.FetchArticleRequestInterval,
		FetchArticleRequestLimit:    apiCfg.ClientConfig.ApiConfig.FetchArticleRequestLimit,
		FetchTrendRequestInterval:   apiCfg.ClientConfig.ApiConfig.FetchTrendRequestInterval,
		FetchTrendRequestLimit:      apiCfg.ClientConfig.ApiConfig.FetchTrendRequestLimit,
	}
	// ClientConfigをData.ClientConfigに変換
	clientConfig := ClientConfig{
		UiConfig:  uiCfg,
		ApiConfig: apiCfgDb,
	}
	// 検索履歴と閲覧記事履歴と購読サイトとお気に入りサイト・記事は変換しない
	return User{
		UserName:     apiCfg.UserName,
		UserUniqueID: apiCfg.UserUniqueID,
		ClientConfig: clientConfig,
	}
}
