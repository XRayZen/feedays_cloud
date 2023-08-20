package Data

type ExploreCategory struct {
	// カテゴリー名
	CategoryName string `json:"categoryName"`
	// カテゴリーの画像URL
	CategoryImage string `json:"categoryImage"`
	// カテゴリーの説明
	CategoryDescription string `json:"categoryDescription"`
	// カテゴリーのID
	CategoryID string `json:"categoryID"`
	// カテゴリー国
	CategoryCountry string `json:"categoryCountry"`
	// カテゴリーのURLはない
	// ユーザーがタップしたら＃カテゴリー名のキーワード検索が走り、その結果を表示する
}
