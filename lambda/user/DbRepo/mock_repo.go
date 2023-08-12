package DBRepo

import "user/Data"

type MockDBRepo struct {
}

func (repo MockDBRepo) ConnectDB(isMock bool) error {
	return nil
}

func (repo MockDBRepo) AutoMigrate() error {
	return nil
}

func (repo MockDBRepo) SearchUserConfig(user_unique_Id string) (Data.UserConfig, error) {
	return Data.UserConfig{}, nil
}

func (repo MockDBRepo) RegisterUser(userInfo Data.UserConfig) error {
	return nil
}

func (repo MockDBRepo) DeleteUser(userId string) error {
	return nil
}

func (repo MockDBRepo) UpdateUser(userId string, configInfo Data.UserConfig) error {
	return nil
}

func (repo MockDBRepo) AddReadHistory(userId string, activityInfo Data.ReadHistory) error {
	return nil
}

// SearchReadHistory(user_unique_Id string, limit int) ([]Data.ReadHistory, error)
func (repo MockDBRepo) SearchReadHistory(userId string, limit int) ([]Data.ReadHistory, error) {
	return []Data.ReadHistory{}, nil
}

func (repo MockDBRepo) ModifySearchHistory(userId string, text string, isAddOrRemove bool) ([]string, error) {
	return []string{}, nil
}

func (repo MockDBRepo) ModifyFavoriteSite(userId string, site_url string, isAddOrRemove bool) error {
	return nil
}

func (repo MockDBRepo) ModifyFavoriteArticle(userId string, articleUrl string, isAddOrRemove bool) error {
	return nil
}

func (repo MockDBRepo) FetchAPIRequestLimit(user_unique_Id string) (Data.ApiConfig, error) {
	return Data.ApiConfig{}, nil
}
