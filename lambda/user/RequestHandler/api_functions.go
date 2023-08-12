package RequestHandler

import (
	"encoding/json"
	"user/DBRepo"
	"user/Data"
)

type APIFunctions struct {
	repo DBRepo.DBRepo
	ip   string
}

func (s APIFunctions) RegisterUser(userId string, userCfgJson string) (string, error) {
	// ユーザー設定をjsonから変換してDBに登録する
	var userConfig Data.UserConfig
	if err := json.Unmarshal([]byte(userCfgJson), &userConfig); err != nil {
		return "", err
	}
	if err := s.repo.RegisterUser(userConfig); err != nil {
		return "", err
	}
	return "Success RegisterUser", nil
}

// 設定を同期する為にユーザー設定を取得する
func (s APIFunctions) ConfigSync(userId string) (string, error) {
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

// サイト・記事閲覧などのアクテビティを記録する
func (s APIFunctions) ReportReadActivity(userId string, readActivityJson string) (string, error) {
	var activityInfo Data.ReadHistory
	if err := json.Unmarshal([]byte(readActivityJson), &activityInfo); err != nil {
		return "", err
	}
	if err := s.repo.AddReadHistory(userId, activityInfo); err != nil {
		return "", err
	}
	// この機能ではAPIアクテビティに記録はしない
	return "Success ReportReadActivity", nil
}

func (s APIFunctions) UpdateConfig(userId string, userCfgJson string) (string, error) {
	// UI設定を変更したらクラウドに送信してクラウドの設定を上書きする
	var userConfig Data.UserConfig
	if err := json.Unmarshal([]byte(userCfgJson), &userConfig); err != nil {
		return "", err
	}
	if err := s.repo.UpdateUser(userId, userConfig); err != nil {
		return "", err
	}
	return "Success UpdateConfig", nil
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
	if err := s.repo.ModifyFavoriteSite(userId, webSite.SiteURL, isAddOrRemove); err != nil {
		return "", err
	}
	return "Success ModifyFavoriteSite", nil
}

func (s APIFunctions) ModifyFavoriteArticle(userId string, articleJson string, isAddOrRemoveBoolJson string) (string, error) {
	var isAddOrRemove bool
	// JSONの文字列をboolに変換する
	if err := json.Unmarshal([]byte(isAddOrRemoveBoolJson), &isAddOrRemove); err != nil {
		return "", err
	}
	var article Data.Article
	if err := json.Unmarshal([]byte(articleJson), &article); err != nil {
		return "", err
	}
	if err := s.repo.ModifyFavoriteArticle(userId, article.Link, isAddOrRemove); err != nil {
		return "", err
	}
	return "Success ModifyFavoriteArticle", nil
}

// APIリクエスト制限を取得して返す 起動時に呼び出される
func (s APIFunctions) GetAPIRequestLimit(userId string) (string, error) {
	// APIリクエスト制限を取得する
	apiCfg, err := s.repo.FetchAPIRequestLimit(userId)
	if err != nil {
		return "", err
	}
	// APIリクエスト制限をjsonに変換する
	apiRequestLimitJson, err := json.Marshal(apiCfg)
	if err != nil {
		return "", err
	}
	return string(apiRequestLimitJson), nil
}

// APIリクエスト制限を更新する
func (s APIFunctions) UpdateAPIRequestLimit(userId string, apiRequestLimitJson string) (string, error) {
	// APIリクエスト制限をjsonから変換する
	var apiCfg Data.ApiConfig
	if err := json.Unmarshal([]byte(apiRequestLimitJson), &apiCfg); err != nil {
		return "", err
	}
	// APIリクエスト制限を更新する
	if err := s.repo.UpdateAPIRequestLimit(userId, apiCfg); err != nil {
		return "", err
	}
	return "Success UpdateAPIRequestLimit", nil
}

// Userと紐付けられている全てのデータを削除する
func (s APIFunctions) DeleteUserData(userId string, isUnscope string) (string, error) {
	if isUnscope == "true" {
		// ユーザーの全てのデータを削除する
		if err := s.repo.DeletesUnscopedUserData(userId); err != nil {
			return "", err
		}
	} else {
		if err := s.repo.DeleteUserData(userId); err != nil {
			return "", err
		}
	}
	return "Success DeleteUserData", nil
}
