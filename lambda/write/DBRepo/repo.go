package DBRepo

import "write/Data"

// DBRepo はDBにアクセスするためのインターフェース
type DBRepo interface {
	GetUserInfo(userId string) (Data.UserInfo, error)
	RegisterUser(userId string, userInfo Data.UserConfig) error
	AddApiActivity(userId string, activityInfo Data.Activity) error
	AddReadActivity(userId string, activityInfo Data.ReadActivity) error
	UpdateConfig(userId string, configInfo Data.UserConfig) error
}

type DBRepoImpl struct {
}

func (repo DBRepoImpl) GetUserInfo(userId string) (Data.UserInfo, error) {
	return Data.UserInfo{}, nil
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
