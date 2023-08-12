package Repo

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

var (
	// AWS公式が推奨するやり方
	// この方がコストが安い
	// https://docs.aws.amazon.com/ja_jp/secretsmanager/latest/userguide/integrating_caching_clientapps.html
	secret_cache, _ = secretcache.New()
)

type secretData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func FetchDbSecret() (secret secretData, err error) {
	// シークレットキャッシュの設定でversionStageを指定する
	secret_cache.VersionStage = os.Getenv("secret_stage")
	// シークレットネームはDBユーザー名と同じにしている
	secret_name := os.Getenv("db_username")
	resJson, err := secret_cache.GetSecretStringWithStage(secret_name, secret_cache.VersionStage)
	if err != nil {
		return secretData{}, err
	}
	// パースしたjsonを構造体に入れる
	var res secretData
	if err := json.Unmarshal([]byte(resJson), &res); err != nil {
		return secretData{}, err
	}
	return res, nil
}
