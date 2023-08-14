package Repo

import (
	"read/Data"
	"time"

	"gorm.io/gorm"
)

// GORMで使う構造型のメンバは大文字で始めないといけない
type User struct {
	gorm.Model
	UserName          string
	UserUniqueID      string `gorm:"uniqueIndex"`
	AccountType       string
	Country           string
	UiConfig          UiConfig
	ReadHistories     []ReadHistory `gorm:"foreignKey:UserID"`
	FavoriteSites     []FavoriteSite
	FavoriteArticles  []FavoriteArticle
	SubscriptionSites []SubscriptionSite
	SearchHistories   []SearchHistory
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
	Link           string
	AccessAt       time.Time
	AccessPlatform string
	AccessIP       string
}

type ApiLimitConfig struct {
	gorm.Model
	AccountType                 string
	RefreshArticleInterval      int
	FetchArticleRequestInterval int
	FetchArticleRequestLimit    int
	FetchTrendRequestInterval   int
	FetchTrendRequestLimit      int
}

type UiConfig struct {
	gorm.Model
	UserID                  uint
	ThemeColorValue         int
	ThemeMode               string
	DrawerMenuOpacity       float64
	MobileTextSize          int
	TabletTextSize          int
	ArticleListMobileSize   float64
	ArticleListTabletSize   float64
	ArticleDetailMobileSize float64
	ArticleDetailTabletSize float64
}

// DB構造体からAPI構造体への変換
func ConvertToApiUserConfig(dbCfg User) (resUserCfg Data.UserConfig) {
	// UiResponsiveFontSizeをData.UiResponsiveFontSizeに変換
	articleListFontSize := Data.UiResponsiveFontSize{
		Mobile:  dbCfg.UiConfig.ArticleListMobileSize,
		Tablet:  dbCfg.UiConfig.ArticleListTabletSize,
		Default: dbCfg.UiConfig.ArticleDetailTabletSize,
	}
	articleDetailFontSize := Data.UiResponsiveFontSize{
		Mobile:  dbCfg.UiConfig.ArticleDetailMobileSize,
		Tablet:  dbCfg.UiConfig.ArticleDetailTabletSize,
		Default: dbCfg.UiConfig.ArticleDetailTabletSize,
	}
	uiCfg := Data.UiConfig{
		ThemeColorValue:       dbCfg.UiConfig.ThemeColorValue,
		ThemeMode:             dbCfg.UiConfig.ThemeMode,
		DrawerMenuOpacity:     dbCfg.UiConfig.DrawerMenuOpacity,
		ArticleListFontSize:   articleListFontSize,
		ArticleDetailFontSize: articleDetailFontSize,
	}
	// 検索履歴をData.SearchHistoryの配列に変換
	var searchHistory []Data.SearchHistory
	for _, searchHistoryDb := range dbCfg.SearchHistories {
		searchHistory = append(searchHistory, Data.SearchHistory{
			SearchWord: searchHistoryDb.SearchWord,
			// 検索履歴の日付をRFC3339に変換
			SearchAt: searchHistoryDb.CreatedAt.Format(time.RFC3339),
		})
	}
	// 読んだ記事をData.ReadHistoryの配列に変換
	var readHistory []Data.ReadHistory
	for _, readHistoryDb := range dbCfg.ReadHistories {
		readHistory = append(readHistory, Data.ReadHistory{
			Link:           readHistoryDb.Link,
			AccessAt:       readHistoryDb.AccessAt.Format(time.RFC3339),
			AccessPlatform: readHistoryDb.AccessPlatform,
			AccessIP:       readHistoryDb.AccessIP,
		})
	}
	// 購読サイトをData.SubscribeWebSiteの配列に変換
	var subscribeWebSite []Data.SubscribeWebSite
	for _, subscribeWebSiteDb := range dbCfg.SubscriptionSites {
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
	for _, favoriteSiteDb := range dbCfg.FavoriteSites {
		favoriteSite = append(favoriteSite, Data.FavoriteSite{
			SiteID:    favoriteSiteDb.SiteID,
			CreatedAt: favoriteSiteDb.CreatedAt.Format(time.RFC3339),
		})
	}
	// お気に入り記事をData.FavoriteArticleの配列に変換
	var favoriteArticle []Data.FavoriteArticle
	for _, favoriteArticleDb := range dbCfg.FavoriteArticles {
		favoriteArticle = append(favoriteArticle, Data.FavoriteArticle{
			ArticleID: favoriteArticleDb.ArticleID,
			CreatedAt: favoriteArticleDb.CreatedAt.Format(time.RFC3339),
		})
	}
	clientConfig := Data.ClientConfig{
		UiConfig: uiCfg,
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

// ユーザー・API設定の更新時に使用する
func ConvertToDbUiCfg(apiCfg Data.UserConfig) (dbUiCfg UiConfig) {
	//フォントサイズをData.UiResponsiveFontSizeからDb.UiResponsiveFontSizeに変換
	uiCfg := UiConfig{
		ThemeColorValue:         apiCfg.ClientConfig.UiConfig.ThemeColorValue,
		ThemeMode:               apiCfg.ClientConfig.UiConfig.ThemeMode,
		DrawerMenuOpacity:       apiCfg.ClientConfig.UiConfig.DrawerMenuOpacity,
		ArticleListMobileSize:   apiCfg.ClientConfig.UiConfig.ArticleListFontSize.Mobile,
		ArticleListTabletSize:   apiCfg.ClientConfig.UiConfig.ArticleListFontSize.Tablet,
		ArticleDetailMobileSize: apiCfg.ClientConfig.UiConfig.ArticleDetailFontSize.Mobile,
		ArticleDetailTabletSize: apiCfg.ClientConfig.UiConfig.ArticleDetailFontSize.Tablet,
	}
	return uiCfg
}

// DataのReadHistoryをDbのReadHistoryに変換
func ConvertToDbReadHistory(readHistory Data.ReadHistory) (ReadHistory, error) {
	// アクセス日時をRFC3339からtime.Timeに変換
	accessAt, err := time.Parse(time.RFC3339, readHistory.AccessAt)
	if err != nil {
		return ReadHistory{}, err
	}
	return ReadHistory{
		Link:           readHistory.Link,
		AccessAt:       accessAt,
		AccessPlatform: readHistory.AccessPlatform,
		AccessIP:       readHistory.AccessIP,
	}, nil
}

// ConvertToApiReadHistory DataのReadHistoryをApiのReadHistoryに変換
func ConvertToApiReadHistory(readHistory ReadHistory, user_unique_id string) Data.ReadHistory {
	return Data.ReadHistory{
		Link:           readHistory.Link,
		AccessAt:       readHistory.AccessAt.Format(time.RFC3339),
		AccessPlatform: readHistory.AccessPlatform,
		AccessIP:       readHistory.AccessIP,
	}
}
