package DBRepo

import (
	"read/Data"
)

type MockDBRepo struct {
}

func (s MockDBRepo) GetUserInfo(userID string) (resUserInfo Data.UserInfo, err error) {
	return Data.UserInfo{}, nil
}

func (s MockDBRepo) GetExploreCategories(userID string, country string) (resExp Data.ExploreCategories, err error) {
	return Data.ExploreCategories{
		CategoryName:        "CategoryName",
	}, nil
}

func (s MockDBRepo) GetRanking(userID string, country string) (resRanking Data.Ranking, err error) {
	return Data.Ranking{
		RankingName:        "RankingName",
	}, nil
}


























