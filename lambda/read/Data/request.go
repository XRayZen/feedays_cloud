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

// 記事取得は最新の記事か最古の記事かを取得するか指定する
// 新規取得は100件上限に取得する
type FetchArticlesRequest struct {
	SiteUrl string `json:"siteUrl"`
	// 新規取得か最新記事取得か最古記事取得かのEnum
	RequestType string `json:"requestType"`
	// 更新間隔（分）
	IntervalMinutes int `json:"intervalMinutes"`
	// クライアント側の記事最新日時
	LastModified string `json:"lastModified"`
	// クライアント側の記事最古日時
	OldestModified string `json:"oldestModified"`
}
