package Repo

import (
	"read/Data"
	"time"
)

type MockDBRepo struct {
}
// readで使う
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


// heavyで使う
func (s MockDBRepo) IsExistSite(site_url string) bool {
	return false
}

func (s MockDBRepo) GetSite(site_url string) (Data.WebSite, error) {
	return Data.WebSite{}, nil
}

func (s MockDBRepo) GetSiteLastModified(site_url string) (time.Time, error) {
	return time.Now(), nil
}

func (s MockDBRepo) RegisterSite(site Data.WebSite, articles []Data.Article) error {
	return nil
}

func (s MockDBRepo) SearchArticlesByKeyword(keyword string) ([]Data.Article, error) {
	return nil, nil
}

func (s MockDBRepo) GetArticlesByTme(siteUrl string, lastModified time.Time) ([]Data.Article, error) {
	return nil, nil
}

func (s MockDBRepo) UpdateArticles(siteUrl string, articles []Data.Article) error {
	return nil
}























