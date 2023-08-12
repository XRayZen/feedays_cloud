package FetchSecret

import (
	// "os"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

var (
	// AWS公式が推奨するやり方
	// この方がコストが安い
	// https://docs.aws.amazon.com/ja_jp/secretsmanager/latest/userguide/integrating_caching_clientapps.html
	secretCache, _ = secretcache.New()
)

type secretData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func FetchDbSecret() (secret secretData, err error) {
	// シークレットキャッシュの設定でversionStageを指定する
	secretCache.VersionStage = os.Getenv("secret_stage")
	// シークレットネームはDBユーザー名と同じにしている
	secret_name := os.Getenv("db_username")
	resJson, err := secretCache.GetSecretStringWithStage(secret_name, secretCache.VersionStage)
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

func DB_Secret() (string, string, error) {
	res,err := FetchDbSecret()
	if err != nil {
		return "", "", err
	}
	return res.Username, res.Password, nil
}

func Secret_read_test() (bool, error) {
	fmt.Println("Secret_read_test Start")
	db_username, db_password, err := DB_Secret()
	if err != nil {
		return false, err
	}
	// ちゃんと取得できているか確認してできたらtrueを返す
	// usernameがadminでpasswordが文字数10以上であることを確認する
	if db_username == "admin" && len(db_password) >= 10 {
		return true, nil
	} else {
		return false, nil
	}
}
