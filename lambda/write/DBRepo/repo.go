package DBRepo

import "write/Data"

// DBRepo はDBにアクセスするためのインターフェース
type DBRepo interface {
	GetUserInfo(userId string) (Data.UserInfo, error)
	GetUserConfig(userId string) (Data.UserConfig, error)
	RegisterUser(userId string, userInfo Data.UserConfig) error
	AddApiActivity(userId string, activityInfo Data.Activity) error
	AddReadActivity(userId string, activityInfo Data.ReadActivity) error
	UpdateConfig(userId string, configInfo Data.UserConfig) error
	// 検索履歴を変更したら履歴を返す
	ModifySearchHistory(userId string, text string, isAddOrRemove bool) ([]string, error)
	ModifyFavoriteSite(userId string, siteInfo Data.WebSite, isAddOrRemove bool) error
	ModifyFavoriteArticle(userId string, articleInfo Data.Article, isAddOrRemove bool) error
	GetAPIRequestLimit(userId string) (Data.ApiRequestLimitConfig, error)
}

type DBRepoImpl struct {
}

func (repo DBRepoImpl) GetUserInfo(userId string) (Data.UserInfo, error) {
	return Data.UserInfo{}, nil
}

func (repo DBRepoImpl) GetUserConfig(userId string) (Data.UserConfig, error) {
	return Data.UserConfig{}, nil
}

func (repo DBRepoImpl) RegisterUser(userId string, userInfo Data.UserConfig) error {
	return nil
}

func (repo DBRepoImpl) UpdateConfig(userId string, configInfo Data.UserConfig) error {
	return nil
}

func (repo DBRepoImpl) AddApiActivity(userId string, activityInfo Data.Activity) error {
	return nil
}

func (repo DBRepoImpl) AddReadActivity(userId string, activityInfo Data.ReadActivity) error {
	return nil
}

func (repo DBRepoImpl) ModifySearchHistory(userId string, text string, isAddOrRemove bool) ([]string, error) {
	return []string{}, nil
}

func (repo DBRepoImpl) ModifyFavoriteSite(userId string, siteInfo Data.WebSite, isAddOrRemove bool) error {
	return nil
}

func (repo DBRepoImpl) ModifyFavoriteArticle(userId string, articleInfo Data.Article, isAddOrRemove bool) error {
	return nil
}

func (repo DBRepoImpl) GetAPIRequestLimit(userId string) (Data.ApiRequestLimitConfig, error) {
	return Data.ApiRequestLimitConfig{}, nil
}