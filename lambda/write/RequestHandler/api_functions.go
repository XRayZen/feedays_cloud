package RequestHandler

import (
	"encoding/json"
	"write/DBRepo"
	"read/Data"
)

type APIFunctions struct {
	repo DBRepo.DBRepo
	ip   string
}



// 設定を同期する
func (s APIFunctions) ConfigSync(userId string, identInfoJson string) (string, error) {
	var identInfo Data.UserAccessIdentInfo
	if err := json.Unmarshal([]byte(identInfoJson), &identInfo); err != nil {
		return "", err
	}
	// ユーザー設定を取得する
	userConfig, err := s.repo.SearchUserConfig(userId)
	if err != nil {
		return "", err
	}
	// レスポンスを返す
	response, err := json.Marshal(Data.ConfigSyncResponse{
		UserConfig:   userConfig,
		ResponseType: "accept",
		Error:        "",
	})
	if err != nil {
		return "", err
	}
	return string(response), nil
}

func (s APIFunctions) RegisterUser(userId string, userCfgJson string, identInfoJson string) (string, error) {
	// ユーザー設定をjsonから変換してDBに登録する
	var userConfig Data.UserConfig
	if err := json.Unmarshal([]byte(userCfgJson), &userConfig); err != nil {
		return "", err
	}
	if err := s.repo.RegisterUser(userConfig); err != nil {
		return "", err
	}
	// アクテビティレコードにユーザー登録イベントを追加する
	identInfo := Data.UserAccessIdentInfo{}
	if err := json.Unmarshal([]byte(identInfoJson), &identInfo); err != nil {
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
	identInfo := Data.UserAccessIdentInfo{}
	if err := json.Unmarshal([]byte(identInfoJson), &identInfo); err != nil {
		return "", err
	}
	if err := s.repo.UpdateUser(userId, userConfig); err != nil {
		return "", err
	}
	return GenAPIResponse("accept", "Success UpdateConfig", "")
}

func (s APIFunctions) ModifySearchHistory(userId string, text string, isAddOrRemoveBool string) (string, error) {
	// Arrayをjsonに変換して返しても良い
	var isAddOrRemove bool
	if isAddOrRemoveBool == "true" {
		isAddOrRemove = true
	} else {
		isAddOrRemove = false
	}
	res_searchHist, err := s.repo.ModifySearchHistory(userId, text, isAddOrRemove)
	if err != nil {
		return "", err
	}
	// resをjsonに変換して返す
	resJson, err := json.Marshal(res_searchHist)
	if err != nil {
		return "", err
	}
	return string(resJson), nil
}

func (s APIFunctions) ModifyFavoriteSite(userId string, webSiteJson string, isAddOrRemoveBool string) (string, error) {
	var isAddOrRemove bool
	if isAddOrRemoveBool == "true" {
		isAddOrRemove = true
	} else {
		isAddOrRemove = false
	}
	var webSite Data.WebSite
	if err := json.Unmarshal([]byte(webSiteJson), &webSite); err != nil {
		return "", err
	}
	if err := s.repo.ModifyFavoriteSite(userId, webSite, isAddOrRemove); err != nil {
		return "", err
	}
	return GenAPIResponse("accept", "Success ModifyFavoriteSite", "")
}

func (s APIFunctions) ModifyFavoriteArticle(userId string, articleJson string, isAddOrRemoveBool string) (string, error) {
	var isAddOrRemove bool
	if isAddOrRemoveBool == "true" {
		isAddOrRemove = true
	} else {
		isAddOrRemove = false
	}
	var article Data.Article
	if err := json.Unmarshal([]byte(articleJson), &article); err != nil {
		return "", err
	}
	if err := s.repo.ModifyFavoriteArticle(userId, article, isAddOrRemove); err != nil {
		return "", err
	}
	return GenAPIResponse("accept", "Success ModifyFavoriteArticle", "")
}

// APIリクエスト制限を取得して返す 起動時に呼び出される
func (s APIFunctions) GetAPIRequestLimit(userId string, identInfoJson string) (string, error) {
	var identInfo Data.UserAccessIdentInfo
	if err := json.Unmarshal([]byte(identInfoJson), &identInfo); err != nil {
		return "", err
	}
	// APIリクエスト制限を取得する
	apiRequestLimit, err := s.repo.GetAPIRequestLimit(userId)
	if err != nil {
		return "", err
	}
	// APIリクエスト制限をjsonに変換する
	apiRequestLimitJson, err := json.Marshal(apiRequestLimit)
	if err != nil {
		return "", err
	}
	return string(apiRequestLimitJson), nil
}
