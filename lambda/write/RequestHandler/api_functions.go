package RequestHandler

import (
	"encoding/json"
	"write/DBRepo"
	"write/Data"
)

type APIFunctions struct {
	repo DBRepo.DBRepo
	ip   string
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

func (s APIFunctions) CodeSync(code string) (string, error) {
	// コード同期はアクティビティを記録しない
	return s.GetUserInfo(code)
}

func (s APIFunctions) RegisterUser(userId string, userCfgJson string, identInfoJson string) (string, error) {
	// ユーザー設定をjsonから変換してDBに登録する
	var userConfig Data.UserConfig
	if err := json.Unmarshal([]byte(userCfgJson), &userConfig); err != nil {
		return "", err
	}
	if err := s.repo.RegisterUser(userId, userConfig); err != nil {
		return "", err
	}
	// アクテビティレコードにユーザー登録イベントを追加する
	identInfo := Data.UserAccessIdentInfo{}
	if err := json.Unmarshal([]byte(identInfoJson), &identInfo); err != nil {
		return "", err
	}
	if err := ReportAPIActivity(s.ip, s.repo, userId, identInfo, "RegisterUser"); err != nil {
		return "", err
	}
	return GenAPIResponse("accept", "Success RegisterUser", "")
}

// サイト・記事閲覧などのアクテビティを記録する
func (s APIFunctions) ReportReadActivity(userId string, readActivity string, identInfo string) (string, error) {
	var activityInfo Data.ReadActivity
	if err := json.Unmarshal([]byte(readActivity), &activityInfo); err != nil {
		return "", err
	}
	if err := s.repo.AddReadActivity(userId, activityInfo); err != nil {
		return "", err
	}
	// この機能ではAPIアクテビティに記録はしない
	return GenAPIResponse("accept", "Success ReportReadActivity", "")
}

func (s APIFunctions) UpdateConfig(userId string, userCfgJson string, identInfoJson string) (string, error) {
	// UI設定を変更したらクラウドに送信してクラウドの設定を上書きする
	var userConfig Data.UserConfig
	if err := json.Unmarshal([]byte(userCfgJson), &userConfig); err != nil {
		return "", err
	}
	var identInfo Data.UserAccessIdentInfo
	if err := json.Unmarshal([]byte(identInfoJson), &identInfoJson); err != nil {
		return "", err
	}
	if err := s.repo.UpdateConfig(userId, userConfig); err != nil {
		return "", err
	}
	if err := ReportAPIActivity(s.ip, s.repo, userId, identInfo, "UpdateConfig"); err != nil {
		return "", err
	}
	return GenAPIResponse("accept", "Success UpdateConfig", "")
}

func (s APIFunctions) ModifySearchHistory(userId string, value string, argument_value2 string) (string, error) {
	
	return "", nil
}

func (s APIFunctions) favoriteSite(userId string, value string, argument_value2 string) (string, error) {
	return "", nil
}

func (s APIFunctions) favoriteArticle(userId string, value string, argument_value2 string) (string, error) {
	return "", nil
}

// APIリクエスト制限を取得して返す
func (s APIFunctions) GetAPIRequestLimit(userId string, identInfoJson string) (string, error) {
	return "", nil
}
