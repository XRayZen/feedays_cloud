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
	// サイトのカテゴリー（ニュース、エンタメ、スポーツ、etc...）
	SiteCategory string `json:"siteCategory"`
	// サイトのRSS URL
	SiteRssURL string `json:"siteRssURL"`
	// サイトのタグ
	SiteTags []string `json:"siteTags"`
	// サイトの最終更新日時
	LastModified string `json:"lastModified"`
}

type Article struct {
	// 記事のタイトル
	Title string `json:"title"`
	// 記事の説明
	Description string `json:"description"`
	// 記事のリンク
	Link string `json:"link"`
	// 記事の画像
	Image RssFeedImage `json:"image"`
	// 記事のサイト名
	Site string `json:"site"`
	// 記事の公開日時
	PublishedAt string `json:"publishedAt"`
	// 記事が既読かどうか
	IsReedLate bool `json:"isReedLate"`
	// 記事のカテゴリー
	Category string `json:"category"`
	// サイトURL
	SiteUrl string `json:"siteUrl"`
}

type RssFeedImage struct {
	// 画像のリンク
	Link string `json:"link"`
	// 画像のデータ
	Image string `json:"image"`
}
