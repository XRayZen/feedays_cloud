package Data

type UserInfo struct {
	// ユーザー名
	UserName string
	// ユーザーID
	UserID string
	// ユーザーがゲストかどうか
	IsGuest bool
	// ユーザーの設定
	Config AppConfig
	// ユーザーのアカウントタイプ(ゲスト, フリー, プロ, アルティメット)
	//
	// 本来enumだが、stringで保存する
	AccountType string
	// ユーザーの検索履歴
	SearchHistory []string
	// ユーザーの国
	UserCountry string
	// ユーザー情報でカテゴリーを持たせる必要がないので削除
	// ユーザーのパスワードも必要ないので削除
}

type UserAccessIdentInfo struct {
	// 端末で取得出来たら
	UUid string `json:"uuid"`
	//accessPlatform
	AccessPlatform string `json:"accessPlatform"`
	//platformType
	PlatformType string `json:"platformType"`
	//brand
	Brand string `json:"brand"`
	//deviceName
	DeviceName string `json:"deviceName"`
	//OSのバージョン
	OsVersion string `json:"osVersion"`
	//isPhysics
	IsPhysics bool `json:"isPhysics"`
}

type ReadActivity struct {
	// ユーザーID
	UserID string `json:"userID"`
	// タイトル
	Title string `json:"title"`
	// リンク
	Link string `json:"link"`
	// アクティビティタイプ (  read Article or Site,)
	Type string `json:"type"`
}
