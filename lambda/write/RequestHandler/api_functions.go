package RequestHandler

import (
	"encoding/json"
	"write/DBRepo"
)

type APIFunctions struct {
	repo DBRepo.DBRepo
}

func (s APIFunctions) GetUserInfo(userId string) (string, error) {
	res, err := s.repo.GetUserInfo(userId)
	if err != nil {
		return "", err
	}
	str, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func (s APIFunctions) CodeSync(code string, value string) (string, error) {
	// ユーザーIDコードを元に同期するアカウントを割り出して設定を読み込んで返す
	return "", nil
}

func (s APIFunctions) RegisterUser(userId string, value string, argument_value2 string) (string, error) {
	return "", nil
}

func (s APIFunctions) ReportActivity(value string) (string, error) {
	return "", nil
}

func (s APIFunctions) SyncConfig(value string) (string, error) {
	// UI設定を変更したらクラウドに送信してクラウドの設定を上書きする
	return "", nil
}

func (s APIFunctions) editRecentSearches(userId string, value string, argument_value2 string) (string, error) {
	return "", nil
}

func (s APIFunctions) favoriteSite(userId string, value string, argument_value2 string) (string, error) {
	return "", nil
}

func (s APIFunctions) favoriteArticle(userId string, value string, argument_value2 string) (string, error) {
	return "", nil
}
