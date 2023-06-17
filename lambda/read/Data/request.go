package Data

type ApiSearchRequest struct {
	// addContent, exploreWeb, powerSearch
	SearchType  string `json:"searchType"`
	Word        string `json:"word"`
	UserID      string `json:"userID"`
	IdentInfo   string `json:"identInfo"`
	AccountType string `json:"accountType"`
	RequestTime string `json:"requestTime"`
}

type FetchFeedRequest struct {
	SiteUrl string `json:"siteUrl"`
	// 更新間隔（分）
	IntervalMinutes int `json:"intervalMinutes"`
	// クライアント側のフィード取得日時
	LastModified string `json:"lastModified"`
}
