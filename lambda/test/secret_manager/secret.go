package SecretManager

import (
	// "os"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

var (
	secretCache, _ = secretcache.New()
)

type secretData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// キャッシュを使ったシークレットバリューの取得
func getSecretValueWithCache(secretName string, secretStage string) (string, error) {
	log.Println("getSecretValueWithCache Start")
	// シークレットキャッシュの設定でversionStageを指定する
	secretCache.VersionStage = os.Getenv("secret_stage")
	result, err := secretCache.GetSecretStringWithStage(secretName, secretStage)
	if err != nil {
		return "", err
	}
	return result, nil
}

func DB_Secret() (string, string, error) {
	// AWS Secrets ManagerからDB接続情報を取得する
	// usernameとpasswordはSecrets Managerに保存している
	fmt.Println("DB_Secret Start")
	region := os.Getenv("region")
	secret_stage := os.Getenv("secret_stage")
	// シークレットネームはDBユーザー名と同じにしている
	secret_name := os.Getenv("db_username")
	log.Println("region:", region)
	log.Println("secret_name:", secret_name)
	log.Println("secret_stage:", secret_stage)
	// ここでシークレットバリューがjsonで帰ってくるのでパースする
	responseJson, err := getSecretValueWithCache(secret_name, secret_stage)
	if err != nil {
		return "", "", err
	}
	log.Println("responseJson:", responseJson)
	// パースしたjsonを構造体に入れる
	var res secretData
	if err := json.Unmarshal([]byte(responseJson), &res); err != nil {
		return "", "", err
	}
	fmt.Println("DB_Secret End")
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
