package DBRepo

import "read/Data"

type MockDBRepo struct {
}

func (repo MockDBRepo) ConnectDB(isMock bool) error {
	return nil
}

func (repo MockDBRepo) AutoMigrate() error {
	return nil
}

func (repo MockDBRepo) GetUserConfig(userId string) (Data.UserConfig, error) {
	return Data.UserConfig{}, nil
}

func (repo MockDBRepo) RegisterUser(userInfo Data.UserConfig) error {
	return nil
}

func (repo MockDBRepo) DeleteUser(userId string) error {
	return nil
}

func (repo MockDBRepo) UpdateAppConfig(userId string, configInfo Data.UserConfig) error {
	return nil
}

func (repo MockDBRepo) AddApiActivity(userId string, activityInfo Data.Activity) error {
	return nil
}

func (repo MockDBRepo) AddReadActivity(userId string, activityInfo Data.ReadActivity) error {
	return nil
}

func (repo MockDBRepo) ModifySearchHistory(userId string, text string, isAddOrRemove bool) ([]string, error) {
	return []string{}, nil
}

func (repo MockDBRepo) ModifyFavoriteSite(userId string, siteInfo Data.WebSite, isAddOrRemove bool) error {
	return nil
}

func (repo MockDBRepo) ModifyFavoriteArticle(userId string, articleInfo Data.Article, isAddOrRemove bool) error {
	return nil
}

func (repo MockDBRepo) GetAPIRequestLimit(userId string) (Data.ApiConfig, error) {
	return Data.ApiConfig{}, nil
}
