package Repo

import (
	"site/Data"
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
	ApiConfig         ApiConfig
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

type ApiConfig struct {
	gorm.Model
	UserID                      uint
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
func ConvertToApiUserConfig(db_user_config User) Data.UserConfig {
	// UiResponsiveFontSizeをData.UiResponsiveFontSizeに変換
	article_list_font_size := Data.UiResponsiveFontSize{
		Mobile:  db_user_config.UiConfig.ArticleListMobileSize,
		Tablet:  db_user_config.UiConfig.ArticleListTabletSize,
		Default: db_user_config.UiConfig.ArticleDetailTabletSize,
	}
	article_detail_font_size := Data.UiResponsiveFontSize{
		Mobile:  db_user_config.UiConfig.ArticleDetailMobileSize,
		Tablet:  db_user_config.UiConfig.ArticleDetailTabletSize,
		Default: db_user_config.UiConfig.ArticleDetailTabletSize,
	}
	ui_config := Data.UiConfig{
		ThemeColorValue:       db_user_config.UiConfig.ThemeColorValue,
		ThemeMode:             db_user_config.UiConfig.ThemeMode,
		DrawerMenuOpacity:     db_user_config.UiConfig.DrawerMenuOpacity,
		ArticleListFontSize:   article_list_font_size,
		ArticleDetailFontSize: article_detail_font_size,
	}
	api_config := Data.ApiConfig{
		RefreshArticleInterval:      db_user_config.ApiConfig.RefreshArticleInterval,
		FetchArticleRequestInterval: db_user_config.ApiConfig.FetchArticleRequestInterval,
		FetchArticleRequestLimit:    db_user_config.ApiConfig.FetchArticleRequestLimit,
		FetchTrendRequestInterval:   db_user_config.ApiConfig.FetchTrendRequestInterval,
		FetchTrendRequestLimit:      db_user_config.ApiConfig.FetchTrendRequestLimit,
	}
	// 検索履歴をData.SearchHistoryの配列に変換
	var search_history []Data.SearchHistory
	for _, search_history_db := range db_user_config.SearchHistories {
		search_history = append(search_history, Data.SearchHistory{
			SearchWord: search_history_db.SearchWord,
			// 検索履歴の日付をRFC3339に変換
			SearchAt: search_history_db.CreatedAt.Format(time.RFC3339),
		})
	}
	// 読んだ記事をData.ReadHistoryの配列に変換
	var read_history []Data.ReadHistory
	for _, read_history_db := range db_user_config.ReadHistories {
		read_history = append(read_history, Data.ReadHistory{
			Link:           read_history_db.Link,
			AccessAt:       read_history_db.AccessAt.Format(time.RFC3339),
			AccessPlatform: read_history_db.AccessPlatform,
			AccessIP:       read_history_db.AccessIP,
		})
	}
	// 購読サイトをData.SubscribeWebSiteの配列に変換
	var subscribe_web_site []Data.SubscribeWebSite
	for _, subscribe_web_site_db := range db_user_config.SubscriptionSites {
		//これだとWebSiteを変換する為に都度クエリが走るから、ここはSiteIDだけ保持してクライアント側で必要になったらSite情報を取得するようにする
		subscribe_web_site = append(subscribe_web_site, Data.SubscribeWebSite{
			SiteID:      subscribe_web_site_db.SiteID,
			FolderIndex: subscribe_web_site_db.FolderIndex,
			FolderName:  subscribe_web_site_db.FolderName,
			SiteIndex:   subscribe_web_site_db.SiteIndex,
		})
	}
	// お気に入りサイトをData.FavoriteSiteの配列に変換
	var favorite_site []Data.FavoriteSite
	for _, favorite_site_db := range db_user_config.FavoriteSites {
		favorite_site = append(favorite_site, Data.FavoriteSite{
			SiteID:    favorite_site_db.SiteID,
			CreatedAt: favorite_site_db.CreatedAt.Format(time.RFC3339),
		})
	}
	// お気に入り記事をData.FavoriteArticleの配列に変換
	var favorite_article []Data.FavoriteArticle
	for _, favorite_article_db := range db_user_config.FavoriteArticles {
		favorite_article = append(favorite_article, Data.FavoriteArticle{
			ArticleID: favorite_article_db.ArticleID,
			CreatedAt: favorite_article_db.CreatedAt.Format(time.RFC3339),
		})
	}
	return Data.UserConfig{
		UserName:     db_user_config.UserName,
		UserUniqueID: db_user_config.UserUniqueID,
		ClientConfig: Data.ClientConfig{
			UiConfig:  ui_config,
			ApiConfig: api_config,
		},
		AccountType:      db_user_config.AccountType,
		SearchHistory:    search_history,
		SubscribeWebSite: subscribe_web_site,
		FavoriteSite:     favorite_site,
		FavoriteArticle:  favorite_article,
		ReadHistory:      read_history,
	}
}

// ユーザー・API設定の更新時に使用する
func ConvertToDbApiCfgAndUiCfg(api_user_config Data.UserConfig) (ApiConfig, UiConfig) {
	//フォントサイズをData.UiResponsiveFontSizeからDb.UiResponsiveFontSizeに変換
	ui_config := UiConfig{
		ThemeColorValue:         api_user_config.ClientConfig.UiConfig.ThemeColorValue,
		ThemeMode:               api_user_config.ClientConfig.UiConfig.ThemeMode,
		DrawerMenuOpacity:       api_user_config.ClientConfig.UiConfig.DrawerMenuOpacity,
		ArticleListMobileSize:   api_user_config.ClientConfig.UiConfig.ArticleListFontSize.Mobile,
		ArticleListTabletSize:   api_user_config.ClientConfig.UiConfig.ArticleListFontSize.Tablet,
		ArticleDetailMobileSize: api_user_config.ClientConfig.UiConfig.ArticleDetailFontSize.Mobile,
		ArticleDetailTabletSize: api_user_config.ClientConfig.UiConfig.ArticleDetailFontSize.Tablet,
	}
	api_config := ApiConfig{
		RefreshArticleInterval:      api_user_config.ClientConfig.ApiConfig.RefreshArticleInterval,
		FetchArticleRequestInterval: api_user_config.ClientConfig.ApiConfig.FetchArticleRequestInterval,
		FetchArticleRequestLimit:    api_user_config.ClientConfig.ApiConfig.FetchArticleRequestLimit,
		FetchTrendRequestInterval:   api_user_config.ClientConfig.ApiConfig.FetchTrendRequestInterval,
		FetchTrendRequestLimit:      api_user_config.ClientConfig.ApiConfig.FetchTrendRequestLimit,
	}
	return api_config, ui_config
}

// DataのReadHistoryをDbのReadHistoryに変換
func ConvertToDbReadHistory(read_history Data.ReadHistory) (ReadHistory, error) {
	// アクセス日時をRFC3339からtime.Timeに変換
	access_at, err := time.Parse(time.RFC3339, read_history.AccessAt)
	if err != nil {
		return ReadHistory{}, err
	}
	return ReadHistory{
		Link:           read_history.Link,
		AccessAt:       access_at,
		AccessPlatform: read_history.AccessPlatform,
		AccessIP:       read_history.AccessIP,
	}, nil
}

func ConvertToApiReadHistory(read_history ReadHistory) Data.ReadHistory {
	return Data.ReadHistory{
		Link:           read_history.Link,
		AccessAt:       read_history.AccessAt.Format(time.RFC3339),
		AccessPlatform: read_history.AccessPlatform,
		AccessIP:       read_history.AccessIP,
	}
}