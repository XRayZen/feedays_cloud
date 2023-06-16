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

type RankingWebSite struct {
	// 順位
	Rank int
	// サイト
	Site WebSite
}

type Ranking struct {
	// ランキング名
	RankingName string
	// ランキングの国
	RankingCountry string
	// ランキングWebサイト
	RankingWebSites []RankingWebSite
}
