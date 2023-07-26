package Data

type UserAccessIdentInfo struct {
	// 端末で取得出来たら
	UUid string `json:"UUID"`
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

type ReadHistory struct {
	// リンク
	Link string `json:"link"`
	// アクセス日時(ISO8601形式-RFC3339)
	AccessAt string `json:"accessAt"`
	// アクセスプラットフォーム
	AccessPlatform string `json:"accessPlatform"`
	// アクセスIP
	AccessIP string `json:"accessIP"`
}
