package Data

type UserInfo struct {
	// ユーザー名
	UserName string `json:"userName"`
	// ユーザーID
	UserID string `json:"userId"`
	// ユーザーがゲストかどうか
	IsGuest bool `json:"isGuest"`
	// ユーザーの設定
	Config AppConfig `json:"config"`
	// ユーザーのアカウントタイプ(ゲスト, フリー, プロ, アルティメット)
	//
	// 本来enumだが、stringで保存する
	AccountType string `json:"accountType"`
	// ユーザーの検索履歴
	SearchHistory []string `json:"searchHistory"`
	// ユーザーの国
	UserCountry string `json:"userCountry"`
	// ユーザー情報でカテゴリーを持たせる必要がないので削除
	// ユーザーのパスワードも必要ないので削除
}

type UserAccessIdentInfo struct {
	// 端末で取得出来たら
	UUid string 			`json:"UUID"`
	//accessPlatform
	AccessPlatform string 	`json:"accessPlatform"`
	//platformType
	PlatformType string 	`json:"platformType"`
	//brand
	Brand string 			`json:"brand"`
	//deviceName
	DeviceName string 		`json:"deviceName"`
	//OSのバージョン
	OsVersion string 		`json:"osVersion"`
	//isPhysics
	IsPhysics bool 			`json:"isPhysics"`
}

type ReadActivity struct {
	// ユーザーID
	UserID string `json:"userID"`
	// タイトル
	Title string `json:"title"`
	// リンク
	Link string `json:"link"`
	// アクティビティタイプ (  Article or Site,)
	Type string `json:"type"`
	// サイト名
	SiteName string `json:"siteName"`
	// サイトURL
	SiteURL string `json:"siteURL"`
	// アクセス日時
	AccessTime string `json:"accessTime"`
	// アクセスプラットフォーム
	AccessPlatform string `json:"accessPlatform"`
	// アクセスIP
	AccessIP string `json:"accessIP"`
}
