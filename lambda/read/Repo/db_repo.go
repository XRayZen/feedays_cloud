package Repo

import (
	"read/Data"
	"time"
)

// テストを容易にするためDBを呼び出す層はインターフェースを定義する
type DBRepository interface {
	GetUserInfo(userID string) (resUserInfo Data.UserInfo, err error)
	GetExploreCategories(userID string, country string) (resExp Data.ExploreCategories, err error)
	GetRanking(userID string, country string) (resRanking Data.Ranking, err error)

	// heavyで使う
	// サイトURLをキーにDBに該当するサイトがあるか確認する
	IsExistSite(site_url string) bool
	// サイトURLをキーにDBに該当するサイトを返す
	FetchSite(site_url string) (Data.WebSite, error)
	// サイトURLをキーに記事更新チェック日時を取得する
	FetchSiteLastModified(site_url string) (time.Time, error)
	// 新規サイトをDB(サイトテーブル)に登録する
	RegisterSite(site Data.WebSite, articles []Data.Article) error
	// キーワード検索でDBに該当する記事を返す
	SearchArticlesByKeyword(keyword string) ([]Data.Article, error)
	// サイトの記事を指定した日時より新しい記事を返す
	GetArticlesByTme(siteUrl string, lastModified time.Time) ([]Data.Article, error)
	// サイトの記事を更新する
	// サイトの更新日時より新しい記事があればDBに登録する
	UpdateArticles(siteUrl string, articles []Data.Article) error
	// サイトを購読登録する
	SubscribeSite(user_id string, siteUrl string, is_subscribe bool) error
}

// DBRepoのリアルを実装
type DBRepoImpl struct {
}

func (s DBRepoImpl) GetUserInfo(userID string) (resUserInfo Data.UserInfo, err error) {
	return Data.UserInfo{}, nil
}

func (s DBRepoImpl) GetExploreCategories(userID string, country string) (resExp Data.ExploreCategories, err error) {
	return Data.ExploreCategories{}, nil
}

func (s DBRepoImpl) GetRanking(userID string, country string) (resRanking Data.Ranking, err error) {
	return Data.Ranking{}, nil
}

// heavyで使う
func (r DBRepoImpl) IsExistSite(site_url string) bool {
	return false
}

func (r DBRepoImpl) FetchSite(site_url string) (Data.WebSite, error) {
	return Data.WebSite{}, nil
}

func (r DBRepoImpl) FetchSiteLastModified(site_url string) (time.Time, error) {
	return time.Now(), nil
}

func (r DBRepoImpl) RegisterSite(site Data.WebSite, articles []Data.Article) error {
	return nil
}

func (r DBRepoImpl) SearchArticlesByKeyword(keyword string) ([]Data.Article, error) {
	return nil, nil
}

func (r DBRepoImpl) GetArticlesByTme(siteUrl string, lastModified time.Time) ([]Data.Article, error) {
	return nil, nil
}

func (r DBRepoImpl) UpdateArticles(siteUrl string, articles []Data.Article) error {
	return nil
}

func (r DBRepoImpl) SubscribeSite(user_id string, siteUrl string, is_subscribe bool) error {
	return nil
}
