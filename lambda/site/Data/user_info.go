package Data

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
