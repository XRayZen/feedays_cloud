package Data

type ClientConfig struct {
	UiConfig  UiConfig  `json:"UiConfig"`
}

type ApiConfig struct {
	AccountType                 string `json:"accountType"`
	RefreshArticleInterval      int    `json:"refreshArticleInterval"`
	FetchArticleRequestInterval int    `json:"fetchArticleRequestInterval"`
	FetchArticleRequestLimit    int    `json:"fetchArticleRequestLimit"`
	FetchTrendRequestInterval   int    `json:"trendRequestInterval"`
	FetchTrendRequestLimit      int    `json:"trendRequestLimit"`
}

type UiConfig struct {
	ThemeColorValue int `json:"themeColorValue"`
	// light, dark, system
	ThemeMode             string               `json:"themeMode"`
	DrawerMenuOpacity     float64              `json:"drawerMenuOpacity"`
	ArticleListFontSize   UiResponsiveFontSize `json:"articleListFontSize"`
	ArticleDetailFontSize UiResponsiveFontSize `json:"articleDetailFontSize"`
}

type UiResponsiveFontSize struct {
	Mobile  float64 `json:"mobile"`
	Tablet  float64 `json:"tablet"`
	Default float64 `json:"defaultSize"`
}
