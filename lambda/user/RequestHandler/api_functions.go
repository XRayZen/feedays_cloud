package RequestHandler

import (
	"encoding/json"
	"log"
	"user/Data"
	"user/DbRepo"
)

type APIFunctions struct {
	db_repo DbRepo.DBRepo
	ip      string
}

func (s APIFunctions) RegisterUser(user_id string, user_config_json string) (string, error) {
	// ユーザー設定をjsonから変換してDBに登録する
	var user_config Data.UserConfig
	if err := json.Unmarshal([]byte(user_config_json), &user_config); err != nil {
		return "", err
	}
	if err := s.db_repo.RegisterUser(user_config); err != nil {
		return "", err
	}
	return "Success RegisterUser", nil
}

// 設定を同期する為にユーザー設定を取得する
func (s APIFunctions) ConfigSync(userId string) (string, error) {
	// ユーザー設定を取得する
	user_config, err := s.db_repo.SearchUserConfig(userId)
	if err != nil {
		return "", err
	}
	// レスポンスを返す
	response, err := json.Marshal(Data.ConfigSyncResponse{
		UserConfig:   user_config,
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
	var read_history Data.ReadHistory
	if err := json.Unmarshal([]byte(readActivityJson), &read_history); err != nil {
		return "", err
	}
	if err := s.db_repo.AddReadHistory(userId, read_history); err != nil {
		return "", err
	}
	// この機能ではAPIアクテビティに記録はしない
	return "Success ReportReadActivity", nil
}

func (s APIFunctions) UpdateUiConfig(userId string, userCfgJson string) (string, error) {
	// UI設定を変更したらクラウドに送信してクラウドの設定を上書きする
	var user_config Data.UserConfig
	if err := json.Unmarshal([]byte(userCfgJson), &user_config); err != nil {
		return "", err
	}
	if err := s.db_repo.UpdateUserUiConfig(userId, user_config); err != nil {
		return "", err
	}
	return "Success UpdateUiConfig", nil
}

func (s APIFunctions) ModifySearchHistory(userId string, text string, is_add_or_remove_bool string) (string, error) {
	// Arrayをjsonに変換して返しても良い
	var is_add_or_remove bool
	if is_add_or_remove_bool == "true" {
		is_add_or_remove = true
	} else {
		is_add_or_remove = false
	}
	// Data.SearchHistoryに変換する
	var history Data.SearchHistory
	if err := json.Unmarshal([]byte(text), &history); err != nil {
		return "", err
	}
	response, err := s.db_repo.ModifySearchHistory(userId, history.SearchWord, is_add_or_remove)
	if err != nil {
		return "", err
	}
	// resをjsonに変換して返す
	response_json, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(response_json), nil
}

func (s APIFunctions) ModifyFavoriteSite(userId string, web_site_json string, is_add_or_remove_bool string) (string, error) {
	var is_add_or_remove bool
	if is_add_or_remove_bool == "true" {
		is_add_or_remove = true
	} else {
		is_add_or_remove = false
	}
	var web_site Data.WebSite
	if err := json.Unmarshal([]byte(web_site_json), &web_site); err != nil {
		return "", err
	}
	if err := s.db_repo.ModifyFavoriteSite(userId, web_site.SiteURL, is_add_or_remove); err != nil {
		return "", err
	}
	return "Success ModifyFavoriteSite", nil
}

func (s APIFunctions) ModifyFavoriteArticle(user_Id string, article_json string, is_add_or_remove_bool_json string) (string, error) {
	var is_add_or_remove bool
	// JSONの文字列をboolに変換する
	if err := json.Unmarshal([]byte(is_add_or_remove_bool_json), &is_add_or_remove); err != nil {
		return "", err
	}
	var article Data.Article
	if err := json.Unmarshal([]byte(article_json), &article); err != nil {
		return "", err
	}
	if err := s.db_repo.ModifyFavoriteArticle(user_Id, article.Link, is_add_or_remove); err != nil {
		return "", err
	}
	return "Success ModifyFavoriteArticle", nil
}

// APIリクエスト制限を取得して返す 起動時に呼び出される
func (s APIFunctions) GetAPIRequestLimit(user_id string) (string, error) {
	// APIリクエスト制限を取得する
	api_config, err := s.db_repo.FetchAPIRequestLimit(user_id)
	if err != nil {
		return "Error GetAPIRequestLimit", err
	}
	// APIリクエスト制限をjsonに変換する
	api_request_limit_json, err := json.Marshal(api_config)
	if err != nil {
		return "", err
	}
	return string(api_request_limit_json), nil
}

// APIリクエスト制限を変更する
func (functions APIFunctions) ModifyAPIRequestLimit(modify_type string, api_request_limit_json string) (string, error) {
	// APIリクエスト制限をjsonから変換する
	var api_config Data.ApiConfig
	if err := json.Unmarshal([]byte(api_request_limit_json), &api_config); err != nil {
		log.Fatalln("Failed ModifyAPIRequestLimit Unmarshal error : ", err)
		return "Error ModifyAPIRequestLimit", err
	}
	// APIリクエスト制限を変更する
	if err := functions.db_repo.ModifyApiRequestLimit(modify_type, api_config); err != nil {
		return "Error ModifyAPIRequestLimit", err
	}
	return "Success ModifyAPIRequestLimit", nil
}

// Userと紐付けられている全てのデータを削除する
func (s APIFunctions) DeleteUserData(user_id string, is_unscope string) (string, error) {
	if is_unscope == "true" {
		// ユーザーの全てのデータを削除する
		if err := s.db_repo.DeletesUnscopedUserData(user_id); err != nil {
			return "", err
		}
	} else {
		if err := s.db_repo.DeleteUserData(user_id); err != nil {
			return "", err
		}
	}
	return "Success DeleteUserData", nil
}
