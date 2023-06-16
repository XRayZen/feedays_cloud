package APIFunction

import (
	"encoding/json"
	"errors"
	"read/Data"
)

func (s APIFunctions) Search(access_ip string, user_id string, request_argument_json1 string) (string, error) {
	// リクエストをjsonから変換する
	var apiSearchRequest = Data.ApiSearchRequest{}
	if err := json.Unmarshal([]byte(request_argument_json1), &apiSearchRequest); err != nil {
		return "", err
	}
	// Wordが空文字の場合はエラー
	if apiSearchRequest.Word == "" {
		return "", errors.New("Word is empty")
	}
	// ワードがURLの場合
	if apiSearchRequest.Word[:4] == "http" {
		// サイトURLをキーにDBに該当するサイトがあるか確認する
	} else {
		// キーワード検索を行う
	}

	return "", nil
}
