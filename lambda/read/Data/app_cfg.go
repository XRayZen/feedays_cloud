package Data

type AppConfig struct {
	ApiRequestConfig ApiRequestLimitConfig `json:"apiRequestConfig"`
	RssFeedConfig    RssFeedConfig         `json:"rssFeedConfig"`
	MobileUiConfig   MobileUiConfig        `json:"mobileUiConfig"`
}

type ApiRequestLimitConfig struct {
	FetchFeedRequestInterval int `json:"fetchFeedRequestInterval"`
	FetchRssFeedRequestLimit int `json:"fetchRssFeedRequestLimit"`
	TrendRequestInterval     int `json:"trendRequestInterval"`
	TrendRequestLimit        int `json:"trendRequestLimit"`
}

type RssFeedConfig struct {
	FeedRefreshInterval int `json:"feedRefreshInterval"`
}

type MobileUiConfig struct {
	ThemeColorValue int `json:"themeColorValue"`
	// light, dark, system
	ThemeMode            string               `json:"themeMode"`
	DrawerMenuOpacity    float64              `json:"drawerMenuOpacity"`
	SiteFeedListFontSize UiResponsiveFontSize `json:"siteFeedListFontSize"`
	FeedDetailFontSize   UiResponsiveFontSize `json:"feedDetailFontSize"`
}

type UiResponsiveFontSize struct {
	Mobile      float64 `json:"mobile"`
	Tablet      float64 `json:"tablet"`
	DefaultSize float64 `json:"defaultSize"`
}
