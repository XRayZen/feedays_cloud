package DBRepo

import (
	"read/Data"
)

// AWS RDBからデータを取得する

// テストを容易にするためDBを呼び出す層はインターフェースを定義する
type DBRepository interface {
	GetUserInfo(userID string) (resUserInfo Data.UserInfo, err error)
	GetExploreCategories(userID string, country string) (resExp Data.ExploreCategories, err error)
	GetRanking(userID string, country string) (resRanking Data.Ranking, err error)
}

// DBRepoのリアルを実装
type DBRepoImpl struct {
}

func (s DBRepoImpl) GetUserInfo(userID string) (resUserInfo Data.UserInfo, err error) {
	return Data.UserInfo{}, nil
}

func (s DBRepoImpl) GetExploreCategories(userID string, country string) (resExp Data.ExploreCategories, err error) {

	return Data.ExploreCategories{
	}, nil
}

func (s DBRepoImpl) GetRanking(userID string, country string) (resRanking Data.Ranking, err error) {
	return Data.Ranking{
	}, nil
}
