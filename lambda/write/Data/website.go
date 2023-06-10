package Data

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