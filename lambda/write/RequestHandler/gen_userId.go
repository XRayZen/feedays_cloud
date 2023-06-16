package RequestHandler

import (
	"encoding/json"
	"write/DBRepo"
	"read/Data"

	// "fmt"
	"math/rand"
)

func GenRandomUserID(repo DBRepo.DBRepo, identInfoJson string, ip string) (string, error) {
	var identInfo Data.UserAccessIdentInfo
	if err := json.Unmarshal([]byte(identInfoJson), &identInfo); err != nil {
		return "", err
	}
	// アクティビティレコードにイベントを追加する
	if err := ReportAPIActivity(ip, repo, "", identInfo, "GenUserID"); err != nil {
		return "", err
	}
	return RandomString(15), nil
}

func RandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
