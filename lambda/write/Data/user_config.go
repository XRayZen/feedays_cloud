package Data

type UserConfig struct {
	UserName      string    `json:"userName"`
	Password      string    `json:"password"`
	UserID        string    `json:"userID"`
	IsGuest       bool      `json:"isGuest"`
	AppConfig     AppConfig `json:"appConfig"`
	AccountType   string    `json:"accountType"`
	SearchHistory []string  `json:"searchHistory"`
	SubscribeWebSite []SubscribeWebSite `json:"subscribeWebSite"`
}

type SubscribeWebSite struct {
	FolderIndex int      `json:"folderIndex"`
	FolderName  string   `json:"folderName"`
	SiteIndex   int      `json:"siteIndex"`
	Site        WebSite  `json:"site"`
}
