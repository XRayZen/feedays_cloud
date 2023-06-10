package Data

import ()

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
	UUid string
	//accessPlatform
	AccessPlatform string
	//platformType
	PlatformType string
	//brand
	Brand string
	//deviceName
	DeviceName string
	//OSのバージョン
	OsVersion string
	//isPhysics
	IsPhysics bool
}

type UserActivity struct {
	// ユーザーID
	UserID string `json:"userID"`
	// IPアドレス
	IP string `json:"ip"`
	// タイトル
	Title string `json:"title"`
	// リンク
	Link string `json:"link"`
	// カテゴリー
	Category string `json:"category"`
	// タグ
	Tags []string `json:"tags"`
	// アクティビティタイプ (  read,subscribe,search,other,)
	Type string `json:"type"`
}
