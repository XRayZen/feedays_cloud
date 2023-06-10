package DBRepo

import "write/Data"

// DBRepo はDBにアクセスするためのインターフェース
type DBRepo interface {
	GenUserID() (string, error)
	GetUserInfo(userId string) (Data.UserInfo, error)
	RegisterUser(userId string, userInfo Data.UserInfo) (bool, error)
	ReportActivity(userId string, activityInfo Data.UserActivity) (bool, error)
	SyncConfig(userId string, configInfo Data.UserInfo) (bool, error)
}

type DBRepoImpl struct {
}

func (repo DBRepoImpl) GenUserID() (string, error) {
	return "", nil
}

func (repo DBRepoImpl) GetUserInfo(userId string) (Data.UserInfo, error) {
	return Data.UserInfo{}, nil
}

func (repo DBRepoImpl) RegisterUser(userId string, userInfo Data.UserInfo) (bool, error) {
	return false, nil
}

func (repo DBRepoImpl) ReportActivity(userId string, activityInfo Data.UserActivity) (bool, error) {
	return false, nil
}

func (repo DBRepoImpl) SyncConfig(userId string, configInfo Data.UserInfo) (bool, error) {
	return false, nil
}
