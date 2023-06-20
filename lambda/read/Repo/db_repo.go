package Repo

import (
	"read/Data"
	"time"
)

// テストを容易にするためDependency Injection（依存性の注入）を採用
// DBを呼び出す層はインターフェースを定義する
type DBRepository interface {
	GetUserInfo(userID string) (resUserInfo Data.UserInfo, err error)
	GetExploreCategories(userID string, country string) (resExp Data.ExploreCategories, err error)

	// heavyで使う
	// サイトURLをキーにサイトテーブルに該当するサイトがあるか確認する
	IsExistSite(site_url string) bool
	// ユーザーはサイトを購読しているか確認する
	IsSubscribeSite(user_id string, site_url string) bool
	// サイトURLをキーにDBに該当するサイトを検索して返す
	FetchSite(site_url string) (Data.WebSite, error)
	// サイトURLをキーに記事更新チェック日時を取得する
	FetchSiteLastModified(site_url string) (time.Time, error)
	// 新規サイトをDB(サイトテーブル)に登録する
	RegisterSite(site Data.WebSite, articles []Data.Article) error
	// キーワード検索でDBに該当する記事を返す
	SearchArticlesByKeyword(keyword string) ([]Data.Article, error)
	// サイトの記事を指定した日時より新しい記事を返す
	SearchArticlesByTime(siteUrl string, lastModified time.Time) ([]Data.Article, error)
	// サイト名をキーにサイトを検索
	SearchSiteByName(siteName string) ([]Data.WebSite, error)
	// サイトの記事を更新する
	// サイトの更新日時より新しい記事があればDBに登録する
	UpdateArticles(siteUrl string, articles []Data.Article) error
	// サイトを購読登録する
	SubscribeSite(user_id string, siteUrl string, is_subscribe bool) error

	// バッチ処理用
	// サイトテーブルを全件取得する
	FetchAllSites() ([]Data.WebSite, error)
	// 閲覧履歴テーブルを全件取得する
	FetchAllHistories() ([]Data.ReadActivity, error)
	// サイトと記事を大量に更新する
	UpdateSitesAndArticles(sites []Data.WebSite, articles []Data.Article) error
	// 時間（From・To）を指定してリードアクテビティを検索する
	SearchReadActivityByTime(from time.Time, to time.Time) ([]Data.ReadActivity, error)
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


// heavyで使う
func (r DBRepoImpl) IsExistSite(site_url string) bool {
	return false
}

func (r DBRepoImpl) IsSubscribeSite(user_id string, site_url string) bool {
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

func (r DBRepoImpl) SearchArticlesByTime(siteUrl string, lastModified time.Time) ([]Data.Article, error) {
	return nil, nil
}

func (r DBRepoImpl) SearchSiteByName(siteName string) ([]Data.WebSite, error) {
	return []Data.WebSite{}, nil
}

func (r DBRepoImpl) UpdateArticles(siteUrl string, articles []Data.Article) error {
	return nil
}

func (r DBRepoImpl) SubscribeSite(user_id string, siteUrl string, is_subscribe bool) error {
	return nil
}

// バッチ処理用
func (r DBRepoImpl) FetchAllSites() ([]Data.WebSite, error) {
	return []Data.WebSite{}, nil
}

func (r DBRepoImpl) FetchAllHistories() ([]Data.ReadActivity, error) {
	return []Data.ReadActivity{}, nil
}

func (r DBRepoImpl) UpdateSitesAndArticles(sites []Data.WebSite, articles []Data.Article) error {
	return nil
}

func (r DBRepoImpl) SearchReadActivityByTime(from time.Time, to time.Time) ([]Data.ReadActivity, error) {
	return []Data.ReadActivity{}, nil
}

