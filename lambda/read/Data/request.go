package Data

type ApiSearchRequest struct {
	// URL or Keyword orSiteName
	// クライアント側でキーワードにURLが入力された場合は検索タイプをURLにしてURL検索を行う
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
