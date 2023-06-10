package Data

type AppConfig struct {
	ApiRequestConfig ApiRequestLimitConfig
	RssFeedConfig    RssFeedConfig
	MobileUiConfig   MobileUiConfig
}

type ApiRequestLimitConfig struct {
	TrendRequestLimit        int
	FetchRssFeedRequestLimit int
	SendActivityMinute       int
}

type RssFeedConfig struct {
	LimitLastFetchTime int
}

type MobileUiConfig struct {
	ThemeColorValue int
	// light, dark, system
	ThemeMode            string
	DrawerMenuOpacity    float64
	SiteFeedListFontSize UiResponsiveFontSize
	FeedDetailFontSize   UiResponsiveFontSize
}

type UiResponsiveFontSize struct {
	Mobile      float64
	Tablet      float64
	DefaultSize float64
}
