package Data

type UserConfig struct {
	UserName         string             `json:"userName"`
	UserUniqueID     string             `json:"userUniqueID"`
	ClientConfig     ClientConfig       `json:"clientConfig"`
	AccountType      string             `json:"accountType"`
	SearchHistory    []SearchHistory    `json:"searchHistory"`
	SubscribeWebSite []SubscribeWebSite `json:"subscribeWebSite"`
	FavoriteSite     []FavoriteSite     `json:"favoriteSite"`
	FavoriteArticle  []FavoriteArticle  `json:"favoriteArticle"`
	ReadHistory      []ReadHistory      `json:"readHistory"`
}

type SubscribeWebSite struct {
	FolderIndex int    `json:"folderIndex"`
	FolderName  string `json:"folderName"`
	SiteIndex   int    `json:"siteIndex"`
	//これだとWebSiteを変換する為に都度クエリが走るから、ここはSiteIDだけ保持してクライアント側で必要になったらSite情報を取得するようにする
	// サイト情報はユーザー側で保持しておくからDB型-API型双変換でサイトクエリが走る事にはならない
	SiteID uint `json:"siteID"`
}

type FavoriteSite struct {
	SiteID    uint   `json:"siteID"`
	CreatedAt string `json:"createdAt"`
}

type FavoriteArticle struct {
	ArticleID uint   `json:"articleID"`
	CreatedAt string `json:"createdAt"`
}

type ReadHistory struct {
	ActivityType   string `json:"activityType"`
	ArticleID      uint   `json:"articleID"`
	SiteID         uint   `json:"siteID"`
	AccessAt       string `json:"accessAt"`
	AccessPlatform string `json:"accessPlatform"`
	AccessIP       string `json:"accessIP"`
}

type SearchHistory struct {
	SearchWord string `json:"searchWord"`
	SearchAt   string `json:"searchAt"`
}
