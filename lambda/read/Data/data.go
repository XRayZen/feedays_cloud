package Data

type UserInfo struct {
	// ユーザーID
	UserID string
	// ユーザー名
	UserName string
	// ユーザーの国
	UserCountry string
	// この関数でDBから取得するのは三種類のみ
}

// WebSiteの構造体
type WebSite struct {
	// サイト名
	SiteName string `json:"siteName"`
	// サイトの画像URL
	SiteImage string `json:"siteImage"`
	// サイトの説明
	SiteDescription string `json:"siteDescription"`
	// サイトのID
	SiteID string `json:"siteID"`
	// サイトのURL
	SiteURL string `json:"siteURL"`
	// サイトのカテゴリー
	SiteCategory string `json:"siteCategory"`
	// サイトのRSS URL
	SiteRssURL string `json:"siteRssURL"`
}

type ExploreCategories struct {
	// カテゴリー名
	CategoryName string `json:"categoryName"`
	// カテゴリーの画像URL
	CategoryImage string `json:"categoryImage"`
	// カテゴリーの説明
	CategoryDescription string `json:"categoryDescription"`
	// カテゴリーのID
	CategoryID string `json:"categoryID"`
	// カテゴリーのURLs
	CategoryURLs []WebSite `json:"categoryURLs"`
}

type RankingWebSite struct {
	// 順位
	Rank int
	// サイト
	Site WebSite
}

type Ranking struct {
	// ランキング名
	RankingName string
	// ランキングの国
	RankingCountry string
	// ランキングWebサイト
	RankingWebSites []RankingWebSite
}





