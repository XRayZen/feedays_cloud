package Data

type UserConfig struct {
	UserName      string    `json:"userName"`
	Password      string    `json:"password"`
	UserID        string    `json:"userID"`
	IsGuest       bool      `json:"isGuest"`
	AppConfig     AppConfig `json:"appConfig"`
	AccountType   string    `json:"accountType"`
	SearchHistory []string  `json:"searchHistory"`
}

type AppConfig struct {
	ApiRequestConfig ApiRequestLimitConfig `json:"apiRequestConfig"`
	RssFeedConfig    RssFeedConfig         `json:"rssFeedConfig"`
	MobileUiConfig   MobileUiConfig        `json:"mobileUiConfig"`
}

type ApiRequestLimitConfig struct {
	TrendRequestLimit        int `json:"trendRequestLimit"`
	FetchRssFeedRequestLimit int `json:"fetchRssFeedRequestLimit"`
	SendActivityMinute       int `json:"sendActivityMinute"`
}

type RssFeedConfig struct {
	LimitLastFetchTime int `json:"limitLastFetchTime"`
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
