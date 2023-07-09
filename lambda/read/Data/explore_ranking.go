package Data

type ExploreCategories struct {
	// カテゴリー名
	CategoryName string `json:"categoryName"`
	// カテゴリーの画像URL
	CategoryImage string `json:"categoryImage"`
	// カテゴリーの説明
	CategoryDescription string `json:"categoryDescription"`
	// カテゴリーのID
	CategoryID string `json:"categoryID"`
	// カテゴリーのURLはない
	// ユーザーがタップしたら＃カテゴリー名のキーワード検索が走り、その結果を表示する
}

type Ranking struct {
	// サイト名
	SiteName string
	Country  string
	// サイトのURL
	SiteURL string
	// サイトの説明
	SiteDescription string
	// サイトの画像URL
	SiteImageURL string
	// サイトのカテゴリー
	SiteCategory string
	// サイトの購読者数
	SubscriberCount int
}
